package disk

import (
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

type Disk struct {
	Total       uint64       `json:"total"`
	Free        uint64       `json:"free"`
	Used        uint64       `json:"used"`
	UsedPercent float64      `json:"used_percent"`
	Partitions  []*Partition `json:"partitions"`
}

type Partition struct {
	Device      string   `json:"device"`
	Mountpoint  string   `json:"mountpoint"`
	Fstype      string   `json:"fstype"`
	Opts        []string `json:"opts"`
	Path        string   `json:"path"`
	Total       uint64   `json:"total"`
	Free        uint64   `json:"free"`
	Used        uint64   `json:"used"`
	UsedPercent float64  `json:"used_percent"`
}

func Stat() (*Disk, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var totalAll uint64
	var usedAll uint64
	var freeAll uint64
	partitions := make([]*Partition, 0)
	for _, p := range parts {
		info, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		if strings.HasPrefix(p.Device, "/dev/loop") {
			continue
		}

		partitions = append(partitions, &Partition{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			Opts:        p.Opts,
			Path:        info.Path,
			Total:       info.Total,
			Free:        info.Free,
			Used:        info.Used,
			UsedPercent: info.UsedPercent,
		})

		totalAll += info.Total
		usedAll += info.Used
		freeAll += info.Free
	}

	var percent float64
	if (usedAll + freeAll) == 0 {
		percent = 0
	} else {
		// We don't use ret.Total to calculate percent.
		// see https://github.com/shirou/gopsutil/issues/562
		percent = (float64(usedAll) / float64(usedAll+freeAll)) * 100.0
	}

	return &Disk{
		Total:       totalAll,
		Free:        usedAll,
		Used:        freeAll,
		UsedPercent: percent,
		Partitions:  partitions,
	}, nil
}
