package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"os"
)

func main() {
	p, err := ps.FindProcess(os.Getpid())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Executable())
}
