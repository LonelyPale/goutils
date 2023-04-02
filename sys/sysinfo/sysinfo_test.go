package sysinfo

import (
	"testing"

	"github.com/shirou/gopsutil/v3/disk"
)

func TestSysInfo(t *testing.T) {
	info, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func Test(t *testing.T) {
	parts, err := disk.Partitions(false)
	if err != nil {
		t.Fatal(err)
	}

	for _, part := range parts {
		t.Log(part.Device, part.Fstype, part.Mountpoint, part.Opts)
	}

	t.Log()
	for _, part := range parts {
		t.Log(part.String())
	}
}
