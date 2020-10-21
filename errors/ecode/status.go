package ecode

import (
	"fmt"
	"strconv"
)

const (
	StatusOK    = 0
	StatusError = 1
)

var statusText = map[int]string{
	StatusOK:    "OK",
	StatusError: "Error",
}

func StatusText(code int) string {
	return statusText[code]
}

// Codes ecode error interface which has a code & message.
type ErrorCode interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
}

type Status struct {
	code    int
	message string
	details []interface{}
}

func New(err error) ErrorCode {
	switch v := err.(type) {
	case ErrorCode:
		return v
	default:
		return Error(StatusError, v.Error())
	}
}

// Error new status with code and message
func Error(code int, message string, objects ...interface{}) *Status {
	return &Status{code, message, append(make([]interface{}, 0), objects...)}
}

// Errorf new status with code and message
func Errorf(code int, format string, args ...interface{}) *Status {
	return Error(code, fmt.Sprintf(format, args...))
}

// Error implement error
func (s *Status) Error() string {
	return s.Message()
}

// Code return error code
func (s *Status) Code() int {
	return s.code
}

// Message return error message for developer
func (s *Status) Message() string {
	if s.message == "" {
		return strconv.Itoa(s.code)
	}
	return s.message
}

// Details return error details
func (s *Status) Details() []interface{} {
	if s == nil {
		return nil
	}
	return s.details
}

// WithDetails WithDetails
func (s *Status) WithDetails(objects ...interface{}) *Status {
	s.details = append(s.details, objects...)
	return s
}
