package sys

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// helper function to make dir creation independent of root dir
func Rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func IsFolderNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// DefaultDataDir is the default data directory
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home == "" {
		return "./"
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support")
	case "windows":
		return filepath.Join(home, "AppData", "Roaming")
	default:
		return home
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
