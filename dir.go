package goutils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// 非空目录
func NonEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	names, _ := f.Readdir(1)
	_ = f.Close()
	return len(names) > 0
}

// file and folder
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// file and folder
func FileNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// 确保完整路径
func EnsureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return fmt.Errorf("could not create directory %v. %v", dir, err)
		}
	}
	return nil
}

// helper function to make dir creation independent of root dir
func Rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// DefaultDataDir is the default data directory
func DefaultDataDir(dirs ...string) string {
	var dir string
	if len(dirs) > 0 {
		dir = dirs[0]
	}

	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home == "" {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", dir)
	case "windows":
		// We used to put everything in %HOME%\AppData\Roaming, but this caused
		// problems with non-typical setups. If this fallback location exists and
		// is non-empty, use it, otherwise DTRT and check %LOCALAPPDATA%.
		fallback := filepath.Join(home, "AppData", "Roaming", dir)
		appdata := WindowsAppData()
		if appdata == "" || NonEmptyDir(fallback) {
			return fallback
		}
		return filepath.Join(appdata, dir)
	default:
		return filepath.Join(home, dir)
	}
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func WindowsAppData() string {
	v := os.Getenv("LOCALAPPDATA")
	if v == "" {
		// Windows XP and below don't have LocalAppData. Crash here because
		// we don't support Windows XP and undefining the variable will cause
		// other issues.
		panic("environment variable LocalAppData is undefined")
	}
	return v
}
