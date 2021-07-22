package goutils

import (
	"os"
	"testing"
)

func TestFileExist(t *testing.T) {
	t.Log(FileExist("/tmp"))
	t.Log(FileExist("/tmp1"))
}

func TestFileNotExist(t *testing.T) {
	t.Log(FileNotExist("/tmp"))
	t.Log(FileNotExist("/tmp1"))
}

func TestGetCurrentPath(t *testing.T) {
	t.Log(GetCurrentPath())

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pwd)
}

func TestHomeDir(t *testing.T) {
	t.Log(HomeDir())
}
