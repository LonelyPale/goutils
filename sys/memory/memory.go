package memory

import "github.com/shirou/gopsutil/v3/mem"

type Memory struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

func Stat() (*Memory, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return &Memory{
		Total:       memory.Total,
		Available:   memory.Available,
		Used:        memory.Used,
		Free:        memory.Free,
		UsedPercent: memory.UsedPercent,
	}, nil
}
