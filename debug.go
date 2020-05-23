package goutils

import (
	"fmt"
	"reflect"
)

func DebugPrintInfo(v reflect.Value) {
	fmt.Println("Kind:", v.Kind())
	fmt.Println("Name:", v.Type().Name())
	fmt.Println("String:", v.Type().String())
	fmt.Println("PkgPath:", v.Type().PkgPath())
	fmt.Println()
}
