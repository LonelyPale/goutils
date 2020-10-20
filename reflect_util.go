package goutils

import (
	"reflect"

	"github.com/LonelyPale/goutils/errors"
)

//去掉指针的包装，以获得原始类型(PrimitiveType)
func PrimitiveValue(i interface{}) interface{} {
	if i == nil {
		panic(errors.ErrCanNotNil)
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
