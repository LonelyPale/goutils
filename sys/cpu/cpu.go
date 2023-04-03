package cpu

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type Cpu struct {
	UsedPercent float64 `json:"used_percent"`
}

func Stat() (*Cpu, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	return &Cpu{UsedPercent: percent[0]}, nil
}
