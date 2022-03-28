package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils/log/rotatelog"
)

// InitLogFile init logrus with hook
func InitLogFile(module string, cfg *Config) error {
	rotateTime, maxAge, err := cfg.Durations()
	if err != nil {
		return err
	}

	if err := clearLockFiles(cfg.LogPath); err != nil {
		return err
	}

	logrus.AddHook(rotatelog.NewRotateHook(cfg.LogPath, module, rotateTime, maxAge))
	logrus.SetOutput(ioutil.Discard) //控制台不输出

	logLevel, err := cfg.Level()
	if err != nil {
		logrus.WithField("error", err).Fatal("wrong log level")
	}

	logrus.SetLevel(logLevel)
	fmt.Printf("all logs are output in the %s directory, log level:%s\n", cfg.LogPath, cfg.LogLevel)
	return nil
}

func clearLockFiles(logPath string) error {
	files, err := ioutil.ReadDir(logPath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, file := range files {
		if ok := strings.HasSuffix(file.Name(), "_lock"); ok {
			if err := os.Remove(filepath.Join(logPath, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
