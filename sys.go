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

var pidFileMap = make(map[string]*FileLock)

func LockPIDFile(pidfile string) error {
	var fileLock *FileLock
	if FileNotExist(pidfile) {
		if err := WritePIDFile(pidfile); err != nil {
			return err
		}

		fileLock = NewFileLock(pidfile)
		if err := fileLock.Lock(); err != nil {
			return err
		}
	} else {
		fileLock = NewFileLock(pidfile)
		if err := fileLock.Lock(); err != nil {
			return err
		}

		if err := WritePIDFile(pidfile); err != nil {
			return err
		}
	}

	pidFileMap[pidfile] = fileLock
	return nil
}

func UnlockPIDFile(pidfile string) error {
	fileLock, ok := pidFileMap[pidfile]
	if !ok {
		return errors.Errorf("file does not exist, %s", pidfile)
	}

	if err := fileLock.Unlock(); err != nil {
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
