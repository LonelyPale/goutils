package main

import (
	"github.com/lonelypale/goutils/cmd"
	"strings"
)

func main() {
	str := "bash -c 'cd /Users/wyb/project/github/chia-blockchain && . ./activate && chia keys generate'"

	for n := 0; n < 100; n++ {
		command, err := cmd.NewCommand(str, cmd.Options{Echo: false})
		if err != nil {
			panic(err)
		}

		res, err := command.CombinedOutput()
		if err != nil {
			panic(err)
		}

		i := strings.LastIndex(res, " ")
		num := res[i+1:]
		if strings.HasPrefix(num, "188") || strings.Index(num, "888") > -1 ||
			strings.HasPrefix(num, "166") || strings.Index(num, "666") > -1 {
			_, _ = cmd.Green.Println(n, "ok:", num)
		} else {
			_, _ = cmd.Blue.Println(n, "pass:", num)
		}
	}

}
