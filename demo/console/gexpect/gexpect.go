package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ThomasRooney/gexpect"
)

func main() {
	//cmd := "sudo pwd"
	cmd := "sudo ls -la"
	pwd := "wyb123456"

	child, err := gexpect.Spawn(cmd)
	if err != nil {
		log.Fatal("Spawn cmd error ", err)
	}

	if err := child.ExpectTimeout("Password:", 3*time.Second); err != nil {
		log.Fatal("Expect timieout error ", err)
	}

	if err := child.SendLine(pwd); err != nil {
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

	go func() {
		if err := child.Wait(); err != nil {
			log.Fatal("Wait error: ", err)
		}
	}()

	fmt.Println("Success")
}
