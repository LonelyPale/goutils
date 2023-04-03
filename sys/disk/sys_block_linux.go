//go:build linux
// +build linux

package disk

import (
	"os"
	"strings"
)

func Devices() ([]string, error) {
	dir, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "loop") {
			continue
		}
		files = append(files, f.Name())
	}

	return files, nil
}
