package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"io"
	"os"
	"path/filepath"

	ebookconvert "QuickPiperAudiobook/internal/binarymanagers/ebookConvert"
	"QuickPiperAudiobook/internal/binarymanagers/iconv"
)

// the function for handling all command logic
func QuickPiperAudiobook(fileName, model string, speakDiacritics bool, outDir string) error {

	rawFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	piper, err := piper.NewPiperClient(model)
	if err != nil {
		return err
	}

	var reader io.Reader
	if !speakDiacritics {
		reader, err = iconv.RemoveDiacritics(reader)
		if err != nil {
			return err
		}
	} else {
		reader = rawFile
	}

	convertedReader, err := ebookconvert.ConvertToText(reader, filepath.Ext(fileName))
	if err != nil {
		return err
	}

	err = piper.Run(fileName, convertedReader, outDir)
	if err != nil {
		return err
	}

	return nil
}
