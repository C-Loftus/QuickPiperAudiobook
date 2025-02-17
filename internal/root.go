package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"QuickPiperAudiobook/internal/lib"
	"QuickPiperAudiobook/internal/parsers/epub"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
}

// make sure the config is not obviously invalid before we try to use it
func sanityCheckConfig(config AudiobookArgs) error {
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
		return fmt.Errorf("currently only epub files can be split into chapters. Please disable the --chapters flag or convert your file to epub")
	}

	return nil
}

func processChapters(piper piper.PiperClient, config AudiobookArgs) (string, error) {
	splitter, err := epub.NewEpubSplitter(config.FileName)
	if err != nil {
		return "", err
	}
	sections, err := splitter.SplitBySection()
	if err != nil {
		return "", err
	}

	errorGroup := errgroup.Group{}
	var mu = &sync.Mutex{}

	// Initialize the slice to store mp3 files in the correct order
	var mp3InOrder = make([]string, len(sections))

	temp_mp3_dir_name, err := os.MkdirTemp("", "piper-ffmpeg-dir-*")
	defer os.RemoveAll(temp_mp3_dir_name)
	if err != nil {
		return "", err
	}

	for i, section := range sections {
		errorGroup.Go(func() error {
			section := section // local variable to capture range variable in local scope
			i := i             // local variable to capture range variable in local scope
			convertedReader, err := ebookconvert.ConvertToText(section.Text, filepath.Ext(section.Filename))
			if err != nil && err != (*ebookconvert.EmptyConversionResultError)(nil) {
				log.Warnf("Warning: Internal file %s was empty when converting %s to a plaintext chapter. Skipping it in the final audiobook. This is ok if it was just images or a titlepage.", section.Filename, config.FileName)
				return nil
			} else if err != nil {
				return err
			}
			if !config.SpeakUTF8 {
				reader, err := iconv.RemoveDiacritics(convertedReader)
				if err != nil {
					return err
				}
				convertedReader = reader
			}

			streamOutput, _, err := piper.Run(section.Filename, convertedReader, config.OutputDirectory, true)
			if err != nil {
				return err
			}

			tmp_mp3_name := filepath.Join(temp_mp3_dir_name, fmt.Sprintf("%d-section-piper-output-%s.mp3", i, section.Filename))

			err = ffmpeg.OutputToMp3(streamOutput.Stdout, tmp_mp3_name)
			if err != nil {
				return err
			}

			// Insert the generated MP3 file inside the list in correct order
			mu.Lock()
			mp3InOrder[i] = tmp_mp3_name
			mu.Unlock()

			return nil
		})
	}
	if err := errorGroup.Wait(); err != nil {
		return "", err
	}

	// filter out empty mp3s which signify chapters with no data
	// i.e. title page or just images
	var filteredMp3InOrder []string
	for _, tmp_mp3_name := range mp3InOrder {
		if tmp_mp3_name == "" {
			continue
		}
		filteredMp3InOrder = append(filteredMp3InOrder, tmp_mp3_name)
	}

	outputName := filepath.Join(config.OutputDirectory, strings.TrimSuffix(filepath.Base(config.FileName), filepath.Ext(config.FileName))) + ".mp3"

	return outputName, ffmpeg.ConcatMp3s(filteredMp3InOrder, outputName)
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

	if err := sanityCheckConfig(config); err != nil {
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

	log.Infof("Audiobook created at: %s", outputName)

	err = beeep.Alert("Audiobook created at "+outputName, "Check the terminal for more info", "")
	if err != nil {
		// although not critical, it's useful to know if the notification failed
		// sometimes a user may not have notify-send in their path
		log.Printf("failed sending alert notification after audiobook completion: %v", err)
	}

	return outputName, nil
}
