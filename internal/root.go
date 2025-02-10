package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"QuickPiperAudiobook/internal/lib"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	ebookconvert "QuickPiperAudiobook/internal/binarymanagers/ebookConvert"
	"QuickPiperAudiobook/internal/binarymanagers/iconv"

	"github.com/gen2brain/beeep"
)

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
}

// Run the core audiobook creation process. Does not include any CLI parsing. Returns the filepath of the created audiobook.
func QuickPiperAudiobook(config AudiobookArgs) (string, error) {
	if config.FileName == "" {
		return "", fmt.Errorf("no file was provided")
	} else if lib.IsUrl(config.FileName) {
		fileNameInUrl := config.FileName[strings.LastIndex(config.FileName, "/")+1:]
		downloadedFile, err := lib.DownloadFile(config.FileName, fileNameInUrl, config.OutputDirectory)
		if err != nil {
			return "", err
		}
		config.FileName = downloadedFile.Name()
	}

	if config.Model == "" {
		return "", fmt.Errorf("no model was provided")
	}
	if config.OutputDirectory == "" {
		return "", fmt.Errorf("no output directory was provided")
	}

	rawFile, err := os.Open(config.FileName)
	if err != nil {
		return "", err
	}

	piper, err := piper.NewPiperClient(config.Model)
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

	var outputName string
	streamOutput, outputFilename, err := piper.Run(config.FileName, convertedReader, config.OutputDirectory, config.OutputAsMp3)
	if err != nil {
		return "", err
	}
	if config.OutputAsMp3 {
		buf := new(bytes.Buffer)
		written, err := io.Copy(buf, streamOutput.Stdout)
		if err != nil {
			return "", fmt.Errorf("failed to read Piper output: %v", err)
		}

		if buf.Len() == 0 || written == 0 {
			return "", fmt.Errorf("piper produced no audio output")
		}

		fileBase := filepath.Base(config.FileName)
		fileNameWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
		outputName = filepath.Join(config.OutputDirectory, fileNameWithoutExt) + ".mp3"

		err = ffmpeg.OutputToMp3(bytes.NewReader(buf.Bytes()), outputName)
		if err != nil {
			return "", err
		}
		fmt.Printf("Audiobook created at: %s\n", outputName)
	} else {
		outputName = outputFilename
		fmt.Printf("Audiobook created at: %s\n", outputFilename)
	}

	err = beeep.Alert("Audiobook created at "+outputName, "Check the terminal for more info", "")
	if err != nil {
		log.Default().Printf("Failed sending notification: %v", err)
	}

	return outputName, nil
}
