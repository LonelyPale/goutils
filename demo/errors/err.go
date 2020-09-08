package main

import (
	"fmt"
	"github.com/LonelyPale/goutils/errors"
)

func main() {
	fmt.Println(test(false))
}

func test(flag bool) (err error) {
	defer func() {
		err = callback(true)
	}()

	if flag {
		return nil
	} else {
		return errors.New("test error")
	}
}

func callback(flag bool) error {
	if flag {
		return nil
	} else {
		return errors.New("callback error")
	}
}
