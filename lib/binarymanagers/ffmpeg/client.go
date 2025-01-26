package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
)

type FfmpegClient struct {
	binary string
}

func NewFFmpegClient() (*FfmpegClient, error) {

	// look for ffmpeg in path
	if _, err := os.Stat("ffmpeg"); os.IsNotExist(err) {
		return &FfmpegClient{binary: "ffmpeg"}, nil
	} else {
		return &FfmpegClient{}, fmt.Errorf("ffmpeg binary not found: %v", err)
	}

}

func (c *FfmpegClient) Run(cmd String) error {
	// construct the command
	cmdList := fmt.Sprintf("%s %s", c.binary, cmd)

	// run the command
	fullCmd := exec.Command("sh", "-c", cmdList)
	output, err := fullCmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %v\n%s", err, string(output))
	}

	return nil
}
