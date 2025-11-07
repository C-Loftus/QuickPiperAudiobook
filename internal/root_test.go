// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/C-Loftus/QuickPiperAudiobook/internal/binarymanagers"

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

	t.Run("end to end with no concurrency; epub has 2 chapters and a title page that is skipped", func(t *testing.T) {

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
			Threads:         1,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))
	})

}

func TestQuickPiperAudiobookWithUTF8(t *testing.T) {

	t.Run("end to end with Chinese", func(t *testing.T) {

		file, err := os.Open(filepath.Join("testdata", "chinese.epub"))
		require.NoError(t, err)
		defer file.Close()

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "zh_CN-huayan-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       true,
			OutputAsMp3:     true,
			Chapters:        true,
			Threads:         3,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".mp3"))

		showChapterCmd := []string{"ffprobe", "-i", outputFilename, "-show_chapters"}
		output, err := binarymanagers.Run(showChapterCmd)
		require.NoError(t, err)
		const heading1InChinese = "标题 1"
		require.Contains(t, output, heading1InChinese)
		const heading2InChinese = "标题 2"
		require.Contains(t, output, heading2InChinese)
	})

	t.Run("end to end with Chinese Markdown", func(t *testing.T) {

		file, err := os.Open(filepath.Join("testdata", "chinese.md"))
		require.NoError(t, err)
		defer file.Close()

		conf := AudiobookArgs{
			FileName:        file.Name(),
			Model:           "zh_CN-huayan-medium.onnx",
			OutputDirectory: ".",
			SpeakUTF8:       true,
			OutputAsMp3:     false,
			Chapters:        true,
			Threads:         4,
		}

		outputFilename, err := QuickPiperAudiobook(conf)
		defer os.Remove(outputFilename)
		require.NoError(t, err)
		_, err = os.Stat(outputFilename)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(outputFilename, ".wav"))

	})
}
