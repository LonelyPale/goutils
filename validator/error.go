package validator

import (
	"bytes"
	"strings"
)

type ValidationErrors []*FieldError

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	var fe *FieldError
	for i := 0; i < len(ve); i++ {
		fe = ve[i]
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

var _ error = new(FieldError)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewFieldError(field string, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

func (fe *FieldError) Error() string {
	return fe.Message
}
