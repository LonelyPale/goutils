package sysinfo

import (
	"fmt"

	"github.com/lonelypale/goutils/sys/cpu"
	"github.com/lonelypale/goutils/sys/disk"
	"github.com/lonelypale/goutils/sys/memory"
)

type SysInfo struct {
	Cpu  *cpu.Cpu       `json:"cpu"`
	Mem  *memory.Memory `json:"mem"`
	Disk *disk.Disk     `json:"disk"`
}

func New() (*SysInfo, error) {
	cpuInfo, err := cpu.Stat()
	if err != nil {
		return nil, fmt.Errorf("SysInfo CPU: %s", err.Error())
	}

	memInfo, err := memory.Stat()
	if err != nil {
		return nil, fmt.Errorf("SysInfo Memory: %s", err.Error())
	}

	diskInfo, err := disk.Stat()
	if err != nil {
		return nil, fmt.Errorf("SysInfo Disk: %s", err.Error())
	}

	return &SysInfo{
		Cpu:  cpuInfo,
		Mem:  memInfo,
		Disk: diskInfo,
	}, nil
}
