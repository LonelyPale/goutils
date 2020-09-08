package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-ps"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pid pid.go
func main() {
	name, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println(os.Getpid(), os.Getppid(), name)
	fmt.Println()

	p, err := ps.FindProcess(os.Getpid())
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if p == nil {
		log.Fatal("should have process")
	}

	if p.Pid() != os.Getpid() {
		log.Fatalf("bad: %#v", p.Pid())
	}

	list, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	for _, proc := range list {
		fmt.Println(proc.Pid(), proc.PPid(), proc.Executable())
	}
}
