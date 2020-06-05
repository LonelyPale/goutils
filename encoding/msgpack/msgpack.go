package msgpack

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func Marshal(v interface{}) ([]byte, error) {
	startTime := time.Now()

	bs, err := msgpack.Marshal(v)
	if err != nil {
		fmt.Println("msgpack Marshal duration:", time.Since(startTime))
		return nil, err
	}

	fmt.Println("msgpack Marshal duration:", time.Since(startTime))
	return bs, nil
}

func Unmarshal(data []byte, v interface{}) error {
	startTime := time.Now()

	err := msgpack.Unmarshal(data, v)
	if err != nil {
		fmt.Println("msgpack Unmarshal duration:", time.Since(startTime))
		return err
	}

	fmt.Println("msgpack Unmarshal duration:", time.Since(startTime))
	return nil
}
