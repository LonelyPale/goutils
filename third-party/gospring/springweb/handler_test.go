package springweb

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	var a string
	typ := reflect.TypeOf(a)
	t.Log(typ.Name())
	t.Log(typ.String())
}

func TestBIND(t *testing.T) {
	t.Log(validBindFn(reflect.TypeOf(demo)))
	t.Log(validBindFn(reflect.TypeOf(demo1)))
	t.Log(validBindFn(reflect.TypeOf(demo2)))
}

func demo() {}

func demo1(a string) {}

func demo2(req *struct {
	name string `json:"n"`
	Age  int    `json:"a"`
}) {
}
