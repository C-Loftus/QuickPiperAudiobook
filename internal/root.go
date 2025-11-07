package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/c-loftus/QuickPiperAudiobook/internal/binarymanagers/piper"
	"github.com/c-loftus/QuickPiperAudiobook/internal/lib"
	"github.com/c-loftus/QuickPiperAudiobook/internal/parsers/epub"

	"github.com/c-loftus/QuickPiperAudiobook/internal/binarymanagers/ffmpeg"

	"github.com/briandowns/spinner"
	log "github.com/charmbracelet/log"

	"github.com/c-loftus/QuickPiperAudiobook/internal/binarymanagers/iconv"

	ebookconvert "github.com/c-loftus/QuickPiperAudiobook/internal/binarymanagers/ebookConvert"

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
		// This is a warning and not an error since we want someone to be able to set chapters = true in the config
		// to use chapters by default for any arbitrary text content and just fall back if it isnt supported
		log.Warnf("Currently only epub files can be split into chapters. Ignoring chapter splitting for %s", config.FileName)
		config.Chapters = false
	}

	if config.Threads > runtime.NumCPU() {
		log.Warnf("%d threads is likely too high for your system; try setting it to a value below %d otherwise may get unexpected I/O errors", config.Threads, runtime.NumCPU())
	}

	return nil
}

// Run the conversion process with chaptered output
// returns the name of the audiobook
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
		log.Warn("Threads value was set to special value 0; ignoring thread limit and using all available resources; this may cause CPU overload")
	} else {
		errorGroup.SetLimit(config.Threads)
	}

	var mu sync.Mutex
	mp3InOrder := make([]ffmpeg.Mp3Section, len(sections))

	tempDir, err := os.MkdirTemp("", "piper-ffmpeg-dir-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	for i, section := range sections {
		i, section := i, section

		errorGroup.Go(func() error {
			section.Filename = strings.ReplaceAll(section.Filename, "/", "_")

			convertedReader, err := ebookconvert.ConvertToText(
				section.Text, filepath.Ext(section.Filename),
			)
			if err != nil {
				var emptyErr *ebookconvert.EmptyConversionResultError
				if errors.As(err, &emptyErr) {
					log.Warnf("Internal file %s was empty when converting and will be skipped. This is expected if it contains just images or no text",
						section.Filename)
					return nil
				}
				return err
			}

			if !config.SpeakUTF8 {
				reader, err := iconv.RemoveDiacritics(convertedReader)
				if err != nil {
					return err
				}
				convertedReader = reader
			}

			// we use a tee reader here so that we can get the title of the audiobook
			// while still being able to pass the full text to piper
			buf := new(bytes.Buffer)
			teeReader := io.TeeReader(convertedReader, buf)
			// 20 is an arbitrary number of bytes to read to get the title
			// the goal is not to have a perfect title but to have something
			// that is reasonably identifiable
			first20 := make([]byte, 20)
			_, err = io.ReadFull(teeReader, first20)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}

			title := strings.TrimSpace(string(first20))

			streamOutput, _, err := piper.Run(section.Filename, io.MultiReader(buf, convertedReader), config.OutputDirectory, true)
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
			log.Debugf("Converted section %d to %s", i, tmpMP3)

			mu.Lock()
			mp3InOrder[i] = ffmpeg.Mp3Section{
				Mp3File: tmpMP3,
				Title:   title,
			}
			mu.Unlock()

			return nil
		})
	}

	if err := errorGroup.Wait(); err != nil {
		return "", err
	}

	var filteredMp3s []ffmpeg.Mp3Section
	for _, section := range mp3InOrder {
		if section.Title != "" {
			filteredMp3s = append(filteredMp3s, section)
		}
	}

	outputName := filepath.Join(
		config.OutputDirectory,
		strings.TrimSuffix(filepath.Base(config.FileName), filepath.Ext(config.FileName))+".mp3",
	)
	log.Debugf("Concatenating %d MP3s", len(filteredMp3s))
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

	start := time.Now()

	config, err := expandHomeDir(config)
	if err != nil {
		return "", err
	}

	if err := sanityCheckConfig(&config); err != nil {
		return "", err
	}

	log.Debugf("Got config after checking and expanding: %+v", config)

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
		log.Errorf("failed sending alert notification after audiobook completion: %v", err)
	}

	elapsed := time.Since(start)
	log.Debugf("Audiobook creation took %.2f seconds (%.2f minutes)", elapsed.Seconds(), elapsed.Minutes())

	return outputName, nil
}
