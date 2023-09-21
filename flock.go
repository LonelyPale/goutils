//go:build !windows
// +build !windows

package goutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/lonelypale/goutils/errors"
)

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

type FileLock struct {
	dir string
	f   *os.File
}

func NewFileLock(dir string) *FileLock {
	return &FileLock{
		dir: dir,
	}
}

func (l *FileLock) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}

	l.f = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}

	return nil
}

func (l *FileLock) Unlock() error {
	defer func() {
		if err := l.f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
