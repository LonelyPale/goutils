package rotatelog

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var defaultFormatter = &logrus.TextFormatter{DisableColors: true}

type rotateHook struct {
	logPath      string
	module       string
	rotationTime time.Duration
	maxAge       time.Duration
	lock         *sync.Mutex
}

func NewRotateHook(logPath, module string, rotationTime, maxAge time.Duration) *rotateHook {
	return &rotateHook{
		lock:         new(sync.Mutex),
		logPath:      logPath,
		module:       module,
		rotationTime: rotationTime,
		maxAge:       maxAge,
	}
}

// Write a log line to an io.Writer.
func (hook *rotateHook) ioWrite(entry *logrus.Entry) error {
	module := hook.module
	if data, ok := entry.Data["module"]; ok {
		module = data.(string)
	}

	logPath := filepath.Join(hook.logPath, module)
	writer, err := rotatelogs.New(
		//logPath+".%Y%m%d%H%M%S",
		logPath+".%Y%m%d.log",

		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logPath+".log"),

		// WithRotationTime设置日志分割的时间，这里设置为一天分割一次
		rotatelogs.WithRotationTime(hook.rotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithRotationCount(maxRemainCnt),
		rotatelogs.WithMaxAge(hook.maxAge),
	)
	if err != nil {
		return err
	}

	msg, err := defaultFormatter.Format(entry)
	if err != nil {
		return err
	}

	if _, err = writer.Write(msg); err != nil {
		return err
	}

	return writer.Close()
}

// Fire write to file
func (hook *rotateHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	return hook.ioWrite(entry)
}

// Levels returns configured log levels.
func (hook *rotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
