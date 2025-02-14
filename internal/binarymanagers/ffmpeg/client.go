package ffmpeg

import (
	"QuickPiperAudiobook/internal/binarymanagers"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Convert raw PCM audio from piper to MP3 using ffmpeg
func OutputToMp3(piperRawAudio io.Reader, outputName string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	cmdStr := fmt.Sprintf("ffmpeg -f s16le -ar 22050 -ac 1 -i pipe:0 -acodec libmp3lame -b:a 128k -y %s", outputName)

	output, err := binarymanagers.RunPiped(cmdStr, piperRawAudio)
	if err != nil {
		return err
	}

	err = output.Handle.Wait()
	if err != nil {
		return err
	}

	validateMp3 := exec.Command("ffmpeg", "-v", "error", "-i", outputName, "-f", "null", "-")

	if err := validateMp3.Run(); err != nil {
		return fmt.Errorf("mp3 output validation failed: %v", err)
	}

	return nil
}

// Concatenate multiple mp3 files into a single mp3
func ConcatMp3s(mp3sInOrder []string, outputName string) error {
	cmdStr := fmt.Sprintf("ffmpeg -i \"concat:%s\" -acodec copy -metadata \"title=Some Song\" %s -map_metadata 0:-1", strings.Join(mp3sInOrder, "|"), outputName)
	_, err := binarymanagers.Run(cmdStr)
	return err
}
