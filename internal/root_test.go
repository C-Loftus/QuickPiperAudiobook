package internal

import (
	ebookconvert "QuickPiperAudiobook/internal/binarymanagers/ebookConvert"
	"QuickPiperAudiobook/internal/binarymanagers/ffmpeg"
	"QuickPiperAudiobook/internal/binarymanagers/iconv"
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"QuickPiperAudiobook/internal/parsers/epub"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/stretchr/testify/require"
)

func TestQuickPiperAudiobookWithWav(t *testing.T) {

	t.Run("end to end with wav and plaintext", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		defer os.Remove(file.Name())
		_, err = file.WriteString("This is some test data that will be converted to speech.")
		require.NoError(t, err)

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     false,
			Chapters:        false,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		err = os.Remove(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".wav"))
	})

	t.Run("end to end with epub and wav", func(t *testing.T) {

		conf := AudiobookArgs{
			FileName:        filepath.Join("testdata", "titlepage_and_2_chapters.epub"),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     false,
			Chapters:        false,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		err = os.Remove(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".wav"))
	})

	t.Run("failure when output dir is non existent", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		defer os.Remove(file.Name())

		const nonexistentDir = "nonexistentDir/foo/bar"
		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: nonexistentDir,
			SpeakUTF8:       false,
			OutputAsMp3:     false,
			Chapters:        false,
		}

		_, err = QuickPiperAudiobook(conf)
		require.Error(t, err)
		require.Contains(t, err.Error(), nonexistentDir)
	})

	t.Run("failure when epub is invalid", func(t *testing.T) {

		const badFile = "testdata/invalid_epub.epub"
		conf := AudiobookArgs{
			FileName:        badFile,
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     false,
			Chapters:        false,
		}

		_, err := QuickPiperAudiobook(conf)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid ZIP file")
	})
}

func TestQuickPiperAudiobookWithMp3(t *testing.T) {

	t.Run("end to end with mp3", func(t *testing.T) {

		file, err := os.CreateTemp("", "*-test.txt")
		require.NoError(t, err)
		defer file.Close()
		_, err = file.WriteString("This is some test data that will be converted to speech.")
		require.NoError(t, err)

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     true,
			Chapters:        false,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))
	})

	t.Run("end to end; epub has one chapter and title page that is skipped", func(t *testing.T) {

		file, err := os.Open(filepath.Join("testdata", "titlepage_and_1_chapter.epub"))
		require.NoError(t, err)
		defer file.Close()

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     true,
			Chapters:        true,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))
	})

	t.Run("end to end; epub has 2 chapters and a title page that is skipped", func(t *testing.T) {

		file, err := os.Open(filepath.Join("testdata", "titlepage_and_2_chapters.epub"))
		require.NoError(t, err)
		defer file.Close()

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "en_US-lessac-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       false,
			OutputAsMp3:     true,
			Chapters:        true,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))
	})

}

// Test the internal functionality of the function that generates an
// epub with chapters
func TestInternalChapterLogic(t *testing.T) {

	config := AudiobookArgs{
		FileName:        filepath.Join("testdata", "titlepage_and_2_chapters.epub"),
		Model:           "en_US-lessac-medium.onnx",
		OutputDirectory: ".",
		SpeakUTF8:       false,
	}

	splitter, err := epub.NewEpubSplitter(config.FileName)
	require.NoError(t, err)
	defer splitter.Close()
	sections, err := splitter.SplitBySection()
	require.NoError(t, err)

	// Initialize the slice to store mp3 files in the correct order
	var mp3InOrder = make([]string, len(sections))

	temp_mp3_dir_name, err := os.MkdirTemp("", "piper-ffmpeg-dir-*")
	defer os.RemoveAll(temp_mp3_dir_name)
	require.NoError(t, err)

	piperClient, err := piper.NewPiperClient(config.Model)
	require.NoError(t, err)

	for i, section := range sections {

		section.Filename = strings.ReplaceAll(section.Filename, "/", "_")

		convertedReader, err := ebookconvert.ConvertToText(section.Text, filepath.Ext(section.Filename))
		if err != nil && err != (*ebookconvert.EmptyConversionResultError)(nil) {
			log.Warnf("Internal epub content '%s' was empty when converting it to a plaintext chapter. Skipping it in the final audiobook. This is ok if it was just images or a cover page.", section.Filename)
			continue
		} else {
			require.NoError(t, err)
		}
		if !config.SpeakUTF8 {
			reader, err := iconv.RemoveDiacritics(convertedReader)
			require.NoError(t, err)
			convertedReader = reader
		}

		streamOutput, _, err := piperClient.Run(section.Filename, convertedReader, config.OutputDirectory, true)
		require.NoError(t, err)

		tmp_mp3_name := filepath.Join(temp_mp3_dir_name, fmt.Sprintf("%d-section-piper-output-%s.mp3", i, section.Filename))

		err = ffmpeg.OutputToMp3(streamOutput.Stdout, tmp_mp3_name)
		require.NoError(t, err)

		// Insert the generated MP3 file inside the list in correct order
		mp3InOrder[i] = tmp_mp3_name

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

	err = ffmpeg.ConcatMp3s(filteredMp3InOrder, outputName)
	require.NoError(t, err)
	defer os.Remove(outputName)

}
