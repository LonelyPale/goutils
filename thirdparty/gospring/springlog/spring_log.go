package springlog

import "github.com/go-spring/spring-logger"

func SetLevel(logLevel string) {
	var level SpringLogger.Level

	switch logLevel {
	case "trace":
		level = SpringLogger.TraceLevel
	case "debug":
		level = SpringLogger.DebugLevel
	case "info":
		level = SpringLogger.InfoLevel
	case "warn":
		level = SpringLogger.WarnLevel
	case "error":
		level = SpringLogger.ErrorLevel
	case "panic":
		level = SpringLogger.PanicLevel
	case "fatal":
		level = SpringLogger.FatalLevel
	default:
		level = SpringLogger.InfoLevel
	}

	SpringLogger.SetLevel(level)
}
