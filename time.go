package goutils

import "time"

var DefaultTimeFormart = "2006-01-02 15:04:05"

func TimeParse(value string) (time.Time, error) {
	return time.Parse(DefaultTimeFormart, value)
}

func TimeParseInLocation(value string) (time.Time, error) {
	return time.ParseInLocation(DefaultTimeFormart, value, time.Local)
}
