//go:build !linux
// +build !linux

package disk

func Devices() ([]string, error) {
	return []string{}, nil
}
