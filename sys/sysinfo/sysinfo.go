package sysinfo

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SysInfo struct {
	Cpu  string `json:"cpu"`
	Mem  string `json:"mem"`
	Disk string `json:"disk"`
}

func New() (*SysInfo, error) {
	cpuInfo, err := CpuPercent()
	if err != nil {
		return nil, fmt.Errorf("SysInfo CPU: %s", err.Error())
	}

	memInfo, err := MemPercent()
	if err != nil {
		return nil, fmt.Errorf("SysInfo Memory: %s", err.Error())
	}

	diskInfo, err := DiskPercent()
	if err != nil {
		return nil, fmt.Errorf("SysInfo Disk: %s", err.Error())
	}

	return &SysInfo{
		Cpu:  fmt.Sprintf("%.2f%%", cpuInfo),
		Mem:  fmt.Sprintf("%.2f%%", memInfo),
		Disk: fmt.Sprintf("%.2f%%", diskInfo),
	}, nil
}

func CpuPercent() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return percent[0], nil
}

func MemPercent() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.UsedPercent, nil
}

func DiskPercent() (float64, error) {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return diskInfo.UsedPercent, nil
}
