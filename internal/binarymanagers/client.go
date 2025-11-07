// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package binarymanagers

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/charmbracelet/log"
)

// Representation of the output of a shell command
type PipedOutput struct {
	Handle *exec.Cmd
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

func RunPiped(cmdName string, args []string, pipedInput io.Reader) (PipedOutput, error) {
	if pipedInput == nil {
		return PipedOutput{}, fmt.Errorf("piped input was nil")
	}

	fullCmd := exec.Command(cmdName, args...)

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
		return PipedOutput{}, fmt.Errorf("command failed when starting: %v", err)
	}

	return PipedOutput{Handle: fullCmd, Stdout: stdout, Stderr: stderr}, nil
}

// Run a shell command and output the combined stdout and stderr
func Run(cmd []string) (string, error) {

	fullCmd := exec.Command(cmd[0], cmd[1:]...)

	outputBytes, err := fullCmd.CombinedOutput()
	if err != nil {
		log.Errorf("Command failed: %v with output: %s", err, string(outputBytes))
		return string(outputBytes), fmt.Errorf("command failed: %v", err)
	}
	return string(outputBytes), nil

}
