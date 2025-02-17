package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"QuickPiperAudiobook/internal/lib"
	"QuickPiperAudiobook/internal/parsers/epub"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/charmbracelet/log"

	ebookconvert "QuickPiperAudiobook/internal/binarymanagers/ebookConvert"
	"QuickPiperAudiobook/internal/binarymanagers/iconv"

	"github.com/gen2brain/beeep"
	"golang.org/x/sync/errgroup"
)

// All the args you can pass to QuickPiperAudiobook
// These are condensed into a struct for easier testing
type AudiobookArgs struct {
	// the file to convert
	FileName string
	// the piper model to use for speech synthesis
	Model string
	// the directory to save the output file
	OutputDirectory string
	// whether to speak utf-8 characters, also known as diacritics
	SpeakUTF8 bool
	// whether to output the audiobook as an mp3 file. if false, use wav
	OutputAsMp3 bool
	// whether to output the audiobook as an mp3 file with chapters
	Chapters bool
	// the number of threads to use when doing concurrent conversions
	Threads int
}

// make sure the config is not obviously invalid before we try to use it
func sanityCheckConfig(config *AudiobookArgs) error {
	if config.FileName == "" {
		return fmt.Errorf("no file was provided")
	}

	if config.Model == "" {
		return fmt.Errorf("no model was provided")
	}

	if config.OutputDirectory == "" {
		return fmt.Errorf("no output directory was provided")
	}

	if _, err := os.Stat(config.OutputDirectory); os.IsNotExist(err) {
		return fmt.Errorf("the output directory %s does not exist", config.OutputDirectory)
	}

	if config.Chapters && filepath.Ext(config.FileName) != ".epub" {
		log.Warnf("currently only epub files can be split into chapters. Ignore chapter splitting for %s", config.FileName)
		config.Chapters = false
	}

	return nil
}

func processChapters(piper piper.PiperClient, config AudiobookArgs) (string, error) {
	splitter, err := epub.NewEpubSplitter(config.FileName)
	if err != nil {
		return "", err
	}
	defer splitter.Close()

	sections, err := splitter.SplitBySection()
	if err != nil {
		return "", err
	}

	errorGroup := errgroup.Group{}

	if config.Threads == 0 {
		log.Warn("threads is set to 0 so will use all available cores; this may cause CPU overload")
	} else {
		errorGroup.SetLimit(config.Threads)
	}

	var mu sync.Mutex

	// Initialize the slice to store mp3 files in the correct order
	mp3InOrder := make([]string, len(sections))

	tempDir, err := os.MkdirTemp("", "piper-ffmpeg-dir-*")
	if err != nil {
		return "", err
	}
	// Clean up temp dir at the end
	defer os.RemoveAll(tempDir)

	for i, section := range sections {
		// capture the loop variables in local scope for the goroutine
		i, section := i, section

		errorGroup.Go(func() error {
			section.Filename = strings.ReplaceAll(section.Filename, "/", "_")

			// Convert to plaintext
			convertedReader, err := ebookconvert.ConvertToText(
				section.Text, filepath.Ext(section.Filename),
			)
			if err != nil {
				var emptyErr *ebookconvert.EmptyConversionResultError
				if errors.As(err, &emptyErr) {
					log.Warnf("File %s was empty (images or cover). Skipping this chapter.",
						section.Filename)
					return nil // skip
				}
				// Otherwise, return the real error so we don't produce a bad MP3
				return err
			}

			if !config.SpeakUTF8 {
				reader, err := iconv.RemoveDiacritics(convertedReader)
				if err != nil {
					return err
				}
				convertedReader = reader
			}

			// Run piper in "streamOutput" mode
			streamOutput, _, err := piper.Run(section.Filename, convertedReader, config.OutputDirectory, true)
			if err != nil {
				return err
			}

			tmpMP3 := filepath.Join(
				tempDir,
				fmt.Sprintf("%04d-section-piper-output-%s.mp3", i, section.Filename),
			)
			err = ffmpeg.OutputToMp3(streamOutput.Stdout, tmpMP3)
			if err != nil {
				return err
			}

			// Place the MP3 path into our slice in the correct index
			mu.Lock()
			mp3InOrder[i] = tmpMP3
			mu.Unlock()

			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := errorGroup.Wait(); err != nil {
		return "", err
	}

	// Filter out empty or skipped chapters
	var filteredMp3s []string
	for _, name := range mp3InOrder {
		if name != "" {
			filteredMp3s = append(filteredMp3s, name)
		}
	}

	// Final output
	outputName := filepath.Join(
		config.OutputDirectory,
		strings.TrimSuffix(filepath.Base(config.FileName), filepath.Ext(config.FileName))+".mp3",
	)

	// Concatenate all final MP3s
	err = ffmpeg.ConcatMp3s(filteredMp3s, outputName)
	if err != nil {
		return "", err
	}

	return outputName, nil
}

// process a book without splitting it into chapters
// returns the filename of the created audiobook
func processWithoutChapters(piper piper.PiperClient, config AudiobookArgs) (string, error) {
	rawFile, err := os.Open(config.FileName)
	if err != nil {
		return "", err
	}

	convertedReader, err := ebookconvert.ConvertToText(rawFile, filepath.Ext(config.FileName))
	if err != nil {
		return "", err
	}

	if !config.SpeakUTF8 {
		reader, err := iconv.RemoveDiacritics(convertedReader)
		if err != nil {
			return "", err
		}
		convertedReader = reader
	}

	streamOutput, piperOutputFilename, err := piper.Run(config.FileName, convertedReader, config.OutputDirectory, config.OutputAsMp3)
	if err != nil {
		return "", err
	}

	var outputName string
	if config.OutputAsMp3 {
		fileBase := filepath.Base(config.FileName)
		fileNameWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
		outputName = filepath.Join(config.OutputDirectory, fileNameWithoutExt) + ".mp3"

		err = ffmpeg.OutputToMp3(streamOutput.Stdout, outputName)
		if err != nil {
			return "", err
		}
	} else {
		outputName = piperOutputFilename
	}

	return outputName, nil

}

// Run the core audiobook creation process. Does not include any CLI parsing. Returns the filepath of the created audiobook.
func QuickPiperAudiobook(config AudiobookArgs) (string, error) {

	config, err := expandHomeDir(config)
	if err != nil {
		return "", err
	}

	if err := sanityCheckConfig(&config); err != nil {
		return "", err
	}

	if lib.IsUrl(config.FileName) {
		fileNameInUrl := config.FileName[strings.LastIndex(config.FileName, "/")+1:]
		downloadedFile, err := lib.DownloadFile(config.FileName, fileNameInUrl, config.OutputDirectory)
		if err != nil {
			return "", err
		}
		config.FileName = downloadedFile.Name()
	}

	piper, err := piper.NewPiperClient(config.Model)
	if err != nil {
		return "", err
	}

	var outputName string
	log.Info("Converting files and generating audiobook. This may take a while...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	if config.Chapters {
		outputName, err = processChapters(*piper, config)
		if err != nil {
			return "", err
		}
	} else {
		outputName, err = processWithoutChapters(*piper, config)
		if err != nil {
			return "", err
		}
	}
	s.Stop()

	log.Infof("Audiobook created at: %s", outputName)

	err = beeep.Alert("Audiobook created at "+outputName, "Check the terminal for more info", "")
	if err != nil {
		// although not critical, it's useful to know if the notification failed
		// sometimes a user may not have notify-send in their path
		log.Warnf("failed sending alert notification after audiobook completion: %v", err)
	}

	return outputName, nil
}
