package log

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Config log config
type Config struct {
	LogPath    string `toml:"log_path"`
	LogLevel   string `toml:"log_level"`
	RotateTime string `toml:"rotate_time"`
	MaxAge     string `toml:"max_age"`
}

var defaultConfig = &Config{
	LogPath:    "logs",
	LogLevel:   "debug",
	RotateTime: "24h",  //每天一个文件
	MaxAge:     "720h", //保留30天
}

// Durations return rotateTime, maxAge time.Duration
func (c *Config) Durations() (rotateTime, maxAge time.Duration, err error) {
	if c.RotateTime == "" {
		c.RotateTime = defaultConfig.RotateTime
	}

	if c.MaxAge == "" {
		c.MaxAge = defaultConfig.MaxAge
	}

	rotateTime, err = time.ParseDuration(c.RotateTime)
	if err != nil {
		return
	}

	maxAge, err = time.ParseDuration(c.MaxAge)
	if err != nil {
		return
	}

	return
}

// Level log level
func (c *Config) Level() (logrus.Level, error) {
	if c.LogLevel == "" {
		c.LogLevel = defaultConfig.LogLevel
	}

	return logrus.ParseLevel(c.LogLevel)
}

func (c *Config) Path() string {
	if c.LogPath == "" {
		c.LogPath = defaultConfig.LogPath
	}

	return c.LogPath
}
