package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A *int
}

func main() {
	t := &T{}
	v := 1
	vptr := &v
	CopyValue(vptr, &t.A) // we pass a reference to t.A since we want to modify it
	fmt.Printf("%v", *t.A)
}

func CopyValue(src interface{}, dest interface{}) {
	srcRef := reflect.ValueOf(src)
	vp := reflect.ValueOf(dest)
	vp.Elem().Set(srcRef)
}
