package goutils

import "testing"

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
}
