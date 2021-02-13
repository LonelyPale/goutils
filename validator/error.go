package validator

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	fieldErrMsg = "%s: %s"
)

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	var fe *fieldError
	for i := 0; i < len(ve); i++ {
		fe = ve[i].(*fieldError)
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

type FieldError interface {
	Error() string
	Field() string
	Label() string
}

// 用来判断 type T 是否实现了接口 I, 用作类型断言, 如果 T 没有实现接口 I, 则编译错误.
var _ FieldError = new(fieldError)
var _ error = new(fieldError)

type fieldError struct {
	field   string
	label   string
	message string
}

func newFieldError(field string, label string, message string) FieldError {
	return &fieldError{
		field:   field,
		label:   label,
		message: message,
	}
}

func (fe *fieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, fe.label, fe.message)
}

func (fe *fieldError) Field() string {
	return fe.field
}

func (fe *fieldError) Label() string {
	return fe.label
}
