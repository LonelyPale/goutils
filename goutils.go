// Created by LonelyPale at 2019-07-27

package goutils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/LonelyPale/goutils/errors"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"
)

// 深拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// 反射注入(struct 结构体)
func Inject(obj interface{}, key string, val interface{}) error {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	} else {
		return errors.ErrMustPointer
	}

	field := value.FieldByName(key)
	if !field.CanSet() {
		// 设置私有字段(小写字段)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	}
	fieldType := field.Type().String()

	v := reflect.ValueOf(val)
	valType := v.Type().String()

	if fieldType == valType {
		field.Set(v)
	} else {
		return errors.Errorf("type mismatch %v to %v", valType, fieldType)
	}

	return nil
}

// Fmt shorthand, XXX DEPRECATED
var Fmt = fmt.Sprintf

// TrapSignal catches the SIGTERM and executes cb function. After that it exits with code 1.
func TrapSignal(cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("captured %v, exiting...\n", sig)
			if cb != nil {
				cb()
			}
			os.Exit(1)
		}
	}()
	select {}
}

func WriteFile(filePath string, contents []byte, mode os.FileMode) error {
	return ioutil.WriteFile(filePath, contents, mode)
}

func MustWriteFile(filePath string, contents []byte, mode os.FileMode) {
	err := WriteFile(filePath, contents, mode)
	if err != nil {
		Exit(Fmt("MustWriteFile failed: %v", err))
	}
}
