package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/google/goexpect"
	"github.com/google/goterm/term"
)

const (
	timeout = 10 * time.Second
)

var (
	passRE   = regexp.MustCompile("password")
	promptRE = regexp.MustCompile("%")
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o test_expect goexpect.go
// 只支持Linux
func main() {
	fmt.Println(term.Bluef("example goexpect"))

	//cmd := "sudo pwd"
	cmd := "sudo ls -la"
	pwd := "wuyb@8btc.com"

	e, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := e.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	s, _, err := e.Expect(passRE, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	if err := e.Send(pwd + "\n"); err != nil {
		log.Fatal(err)
	}

	result, _, err := e.Expect(promptRE, timeout)
	if err != nil {
		if err.Error() != "expect: Process not running" {
			log.Fatal(11, err)
		}
	}
	fmt.Println(term.Green(result))
}
