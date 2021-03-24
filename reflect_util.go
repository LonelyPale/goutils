package goutils

import (
	"reflect"

	"github.com/LonelyPale/goutils/errors"
)

//去掉指针的包装，以获得原始类型的值(PrimitiveType)
func PrimitiveValue(i interface{}) interface{} {
	if i == nil {
		return nil
	}

	val := reflect.ValueOf(i)

	for {
		if k := val.Kind(); k == reflect.Ptr {
			val = val.Elem()
		} else {
			break
		}
	}

	if !val.CanInterface() {
		panic(errors.New("not CanInterface()"))
	}

	return val.Interface()
}

// Indirect 解除 Type 的指针
func Indirect(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}
	return t.Elem()
}
