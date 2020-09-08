package main

import (
	"fmt"
	"regexp"
)

func main() {
	promptRE := regexp.MustCompile("%")
	promptRE1 := regexp.MustCompile(".*")
	fmt.Println(promptRE.FindString("/tmp/123/sdfasdf\n"))
	fmt.Println(promptRE1.FindString("/tmp/123/sdfasdf\n"))
}
