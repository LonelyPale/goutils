package goutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mitchellh/go-ps"

	"github.com/LonelyPale/goutils/errors"
)

// exit app
func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}

func LockPIDFile(pidfile string) error {
	if FileNotExist(pidfile) {
		if err := WritePIDFile(pidfile); err != nil {
			return err
		}

		if err := NewFileLock(pidfile).Lock(); err != nil {
			return err
		}
	} else {
		if err := NewFileLock(pidfile).Lock(); err != nil {
			return err
		}

		if err := WritePIDFile(pidfile); err != nil {
			return err
		}
	}

	return nil
}

func UnlockPIDFile(pidfile string) error {
	if FileNotExist(pidfile) {
		return errors.Errorf("file does not exist, %s", pidfile)
	}

	if err := NewFileLock(pidfile).Unlock(); err != nil {
		return err
	}

	if err := os.Remove(pidfile); err != nil {
		return err
	}

	return nil
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
