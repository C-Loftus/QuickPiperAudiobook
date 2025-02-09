package binarymanagers

import (
	"fmt"
	"io"
	"os/exec"
)

type PipedOutput struct {
	Handle *exec.Cmd
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

func RunPiped(cmd string, pipedInput io.Reader) (PipedOutput, error) {
	if pipedInput == nil {
		return PipedOutput{}, fmt.Errorf("piped input was nil")
	}

	fullCmd := exec.Command("sh", "-c", cmd)

	stdout, err := fullCmd.StdoutPipe()
	if err != nil {
		return PipedOutput{}, fmt.Errorf("failed getting stdout: %v", err)
	}

	stderr, err := fullCmd.StderrPipe()
	if err != nil {
		return PipedOutput{}, fmt.Errorf("failed getting stderr: %v", err)
	}

	fullCmd.Stdin = pipedInput

	if err := fullCmd.Start(); err != nil {
		return PipedOutput{}, fmt.Errorf("command failed: %v", err)
	}

	return PipedOutput{Handle: fullCmd, Stdout: stdout, Stderr: stderr}, nil
}

func Run(cmd string) (string, error) {

	fullCmd := exec.Command("sh", "-c", cmd)

	outputBytes, err := fullCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v", err)
	}
	return string(outputBytes), nil

}
