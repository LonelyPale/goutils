package goutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mitchellh/go-ps"
)

// exit app
func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}

func WritePIDFile(filepath string) error {
	if len(filepath) == 0 {
		if name, err := CurrentProcessName(); err != nil {
			filepath = "temp.pid"
		} else {
			filepath = name + ".pid"
		}
	}

	pid := os.Getpid()
	data := []byte(strconv.Itoa(pid))
	if err := ioutil.WriteFile(filepath, data, 0666); err != nil {
		return err
	}

	return nil
}

func CurrentProcessName() (string, error) {
	pid := os.Getpid()
	process, err := ps.FindProcess(pid)
	if err != nil {
		return "", err
	}

	return process.Executable(), nil
}
