package ref

import (
	"reflect"

	"github.com/LonelyPale/goutils/errors"
)

//去掉指针的包装，以获得原始类型的值
func Primitive(i interface{}) interface{} {
	if i == nil {
		return nil
	}

	val := PrimitiveValue(reflect.ValueOf(i))
	if !val.CanInterface() {
		panic(errors.New("not CanInterface()"))
	}

	return val.Interface()
}

func PrimitiveValue(v reflect.Value) reflect.Value {
	for {
		if k := v.Kind(); k == reflect.Ptr {
			v = v.Elem()
		} else {
			return v
		}
	}
}

func PrimitiveType(t reflect.Type) reflect.Type {
	for {
		if k := t.Kind(); k == reflect.Ptr {
			t = t.Elem()
		} else {
			return t
		}
	}
}

// Indirect 解除 Type 的指针
func Indirect(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}
	return t.Elem()
}

// Indirect 解除 Value 的指针
// reflect.Indirect()
