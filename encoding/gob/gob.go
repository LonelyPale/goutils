package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

func Serialize(value interface{}) ([]byte, error) {
	startTime := time.Now()

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	bs := buf.Bytes()
	fmt.Println("serialize duration:", time.Since(startTime))
	return bs, nil
}

func Deserialize(valueBytes []byte) (interface{}, error) {
	startTime := time.Now()

	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	fmt.Println("deserialize duration:", time.Since(startTime))
	return value, nil
}
