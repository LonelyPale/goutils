package goutils

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-ps"
)

// exit app
func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}

func CurrentProcessName() (string, error) {
	pid := os.Getpid()
	process, err := ps.FindProcess(pid)
	if err != nil {
		return "", err
	}

	return process.Executable(), nil
}
