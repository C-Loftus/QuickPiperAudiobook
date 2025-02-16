package ffmpeg

import (
	"QuickPiperAudiobook/internal/binarymanagers"
	"fmt"
	"io"
	"os/exec"
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

	verifyCmd := fmt.Sprintf("ffmpeg -v error -i %s -f null -", outputName)
	verificationOutput, err := binarymanagers.Run(verifyCmd)

	if err != nil {
		return fmt.Errorf("%v: %s", err, verificationOutput)
	}

	return nil
}
