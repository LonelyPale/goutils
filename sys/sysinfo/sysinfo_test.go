package sysinfo

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

func TestSysInfo(t *testing.T) {
	info, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func TestCpu(t *testing.T) {
	physicalcpu, err := cpu.Counts(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(physicalcpu)

	logicalcpu, err := cpu.Counts(true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(logicalcpu)

	a, _ := cpu.Percent(time.Second, false)
	t.Log(a)
}

func TestMemory(t *testing.T) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(memory)
}

func TestDisk(t *testing.T) {
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
