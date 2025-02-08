package internal

import (
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
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

// the function for handling all command logic
func QuickPiperAudiobook(fileName, model, outDir string, speakDiacritics, outputMp3 bool) error {
	if fileName == "" {
		return fmt.Errorf("no file was provided")
	}
	if model == "" {
		return fmt.Errorf("no model was provided")
	}
	if outDir == "" {
		return fmt.Errorf("no output directory was provided")
	}

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
		reader, err = iconv.RemoveDiacritics(rawFile)
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

	streamOutput, err := piper.Run(fileName, convertedReader, outDir, outputMp3)
	if err != nil {
		return err
	}
	if outputMp3 {
		buf := new(bytes.Buffer)
		written, err := io.Copy(buf, streamOutput.Stdout)
		if err != nil {
			return fmt.Errorf("failed to read Piper output: %v", err)
		}

		if buf.Len() == 0 || written == 0 {
			return fmt.Errorf("piper produced no audio output")
		}

		fileBase := filepath.Base(fileName)
		fileNameWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
		outputName := filepath.Join(outDir, fileNameWithoutExt) + ".mp3"

		err = ffmpeg.OutputToMp3(bytes.NewReader(buf.Bytes()), outputName)
		if err != nil {
			return err
		}
	}

	err = beeep.Alert("Audiobook created", "Check the terminal for more info", "")
	if err != nil {
		log.Default().Printf("Failed sending notification: %v", err)
	}

	return nil
}
