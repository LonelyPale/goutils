//go:build windows
// +build windows

package goutils

func LockPIDFile(pidfile string) error {
	return nil
}

func UnlockPIDFile(pidfile string) error {
	return nil
}
