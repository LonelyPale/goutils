package goutils

import (
	"fmt"
	"os"
)

// exit app
func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}
