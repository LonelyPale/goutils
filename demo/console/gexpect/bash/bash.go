package main

import (
	"fmt"
	"log"

	"github.com/ThomasRooney/gexpect"
)

func main() {
	bash := "bash"
	cmd := "ls -la"

	child, err := gexpect.Spawn(bash)
	if err != nil {
		log.Fatal("Spawn cmd error ", err)
	}

	if err := child.SendLine(cmd); err != nil {
		log.Fatal("SendLine password error ", err)
	}

	for {
		s, err := child.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				log.Fatal(err)
				return
			}
		}
		fmt.Println(s)
	}

	if err := child.SendLine("exit"); err != nil {
		log.Fatal("SendLine cmd error:", err)
	}

	go func() {
		if err := child.Wait(); err != nil {
			log.Fatal("Wait error: ", err)
		}
	}()

	fmt.Println("Success")
}
