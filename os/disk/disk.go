package disk

import (
	"github.com/LonelyPale/goutils/errors"
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	_  = iota
	KB = 1 << (10 * iota) // 2 ^ 10
	MB                    // 2 ^ 20
	GB                    // 2 ^ 30
	TB                    // 2 ^ 40
	PB                    // 2 ^ 50
	EB                    // 2 ^ 60
	ZB                    // 2 ^ 70 //constant 1180591620717411303424 overflows int
	YB                    // 2 ^ 80 //constant 1208925819614629174706176 overflows int
)

func Free(path string) (uint64, error) {
	if path == "" {
		return 0, errors.New("path cannot be empty")
	}

	info, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	//free := info.Total - info.Used
	return info.Free, nil
}
