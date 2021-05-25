package goutils

import "testing"

func TestDiskSize_String(t *testing.T) {
	t.Log(DiskSize(56832).String())
	t.Log(DiskSize(224256).String())
	t.Log(DiskSize(449070816).String())
}
