// Created by LonelyPale at 2019-07-27

package goutils

import (
	"bytes"
	"encoding/gob"
)

// 深拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
