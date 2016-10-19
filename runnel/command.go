package runnel

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"os/exec"
)

// RunCommand ...
func (run *Runnel) RunCommand(executable string,
	args []string) (string, error) {
	key := uuid.NewV4()

	// Redirect command output to buffer
	buf, err := run.Buffer(executable, args)
	if err != nil {
		return "", err
	}

	// Save buffer to redis
	go run.BufferSave(key.String(), buf)

	return key.String(), nil
}

// Buffer will hold command output
func (run *Runnel) Buffer(executable string,
	args []string) (*bufio.Scanner, error) {

	if executable == "" {
		return nil, errors.New("Command executable required")
	}

	command := exec.Command(executable, args...)

	// Get stdout and stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}

	// Start command
	if err := command.Start(); err != nil {
		return nil, err
	}

	// Redirect stdout and stderr to buffer
	multi := io.MultiReader(stdout, stderr)
	return bufio.NewScanner(multi), nil
}

// BufferSave to redis
func (run *Runnel) BufferSave(key string, buffer *bufio.Scanner) string {

	for buffer.Scan() {
		// fmt.Printf("%s\n", in.Text()) // write each line to your log, or anything you need
		if err := run.Redis.LPush(key, buffer.Text()).Err(); err != nil {
			fmt.Println(err)
		}
	}

	if err := buffer.Err(); err != nil {
		// fmt.Printf("error: %s", err)

		if err := run.Redis.LPush(key, err).Err(); err != nil {
			fmt.Println(err)
		}
	}

	return ""
}
