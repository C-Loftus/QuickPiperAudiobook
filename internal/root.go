package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"QuickPiperAudiobook/internal/lib"
	"QuickPiperAudiobook/internal/parsers/epub"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	SpeakDiacritics bool
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

	if config.Chapters && filepath.Ext(config.FileName) != ".epub" {
		return fmt.Errorf("currently only epub files can be split into chapters. Please disable the --chapters flag or convert your file to epub")
	}

	return nil
}

// process a book and split it into chapters
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
	var mp3Files []string
	var mu = &sync.Mutex{}

	for _, section := range sections {
		errorGroup.Go(func() error {
			convertedReader, err := ebookconvert.ConvertToText(section.Text, filepath.Ext(config.FileName))
			if err != nil {
				return err
			}
			if !config.SpeakDiacritics {
				reader, err := iconv.RemoveDiacritics(convertedReader)
				if err != nil {
					return err
				}
				convertedReader = reader
			}

			streamOutput, outputFilename, err := piper.Run(section.Filename, convertedReader, config.OutputDirectory, true)
			if err != nil {
				return err
			}
			fileBase := filepath.Base(outputFilename)
			fileNameWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
			outputName := filepath.Join(config.OutputDirectory, fileNameWithoutExt) + ".mp3"

			mu.Lock()
			mp3Files = append(mp3Files, outputName)
			mu.Unlock()

			err = ffmpeg.OutputToMp3(streamOutput.Stdout, outputName)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := errorGroup.Wait(); err != nil {
		return "", err
	}
	outputName := "test.mp3"
	return outputName, ffmpeg.ConcatMp3s(mp3Files, outputName)
}

// process a book without splitting it into chapters
func processWithoutChapters(piper piper.PiperClient, config AudiobookArgs) (string, error) {
	rawFile, err := os.Open(config.FileName)
	if err != nil {
		return "", err
	}

	var reader io.Reader
	if !config.SpeakDiacritics {
		reader, err = iconv.RemoveDiacritics(rawFile)
		if err != nil {
			return "", err
		}
	} else {
		reader = rawFile
	}

	convertedReader, err := ebookconvert.ConvertToText(reader, filepath.Ext(config.FileName))
	if err != nil {
		return "", err
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

	fmt.Printf("Audiobook created at: %s\n", outputName)

	err = beeep.Alert("Audiobook created at "+outputName, "Check the terminal for more info", "")
	if err != nil {
		// although not critical, it's useful to know if the notification failed
		// sometimes a user may not have notify-send in their path
		fmt.Printf("failed sending alert notification after audiobook completion: %v", err)
	}

	return outputName, nil
}
