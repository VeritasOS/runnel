package runnel

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"os/exec"
)

// RunCommand ...
func (run *Runnel) RunCommand(executable string,
	args []string, cwd string) (string, error) {
	key := uuid.NewV4()

	// Redirect command output to buffer
	buf, err := run.Buffer(executable, args, cwd)
	if err != nil {
		return "", err
	}

	// Save buffer to redis
	go run.BufferSave(key.String(), buf)

	return key.String(), nil
}

// Buffer will hold command output
func (run *Runnel) Buffer(executable string,
	args []string, cwd string) (*bufio.Scanner, error) {

	if executable == "" {
		return nil, errors.New("Command executable required")
	}

	command := exec.Command(executable, args...)

	if cwd != "" {
		command.Dir = cwd
	}
	log.Printf("Executing command %s:%s %s", command.Dir, executable, args)

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

	log.Println("Pushing to redis store", key)

	for buffer.Scan() {
		log.Println(buffer.Text())
		if err := run.Redis.LPush(key, buffer.Text()).Err(); err != nil {
			fmt.Println(err)
		}
	}

	if err := buffer.Err(); err != nil {
		log.Println(err)
		if err := run.Redis.LPush(key, err).Err(); err != nil {
			fmt.Println(err)
		}
	}

	return ""
}
