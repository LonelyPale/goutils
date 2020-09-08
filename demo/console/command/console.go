package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	if err := execCmd("sudo", []string{"ls"}); err != nil {
		fmt.Println(err)
	}
}

func git() {
	if err := execCmd("git", []string{"pull"}); err != nil {
		fmt.Println(err)
	}
}

func ps() {
	if err := execCmd("ps", []string{"-a"}); err != nil {
		fmt.Println(err)
	}
}

func execCmd(shell string, raw []string) error {
	cmd := exec.Command(shell, raw...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		text := s.Text()
		fmt.Println(text)

		if _, err := stdin.Write([]byte("wyb123456")); err != nil {
			fmt.Println(err)
		}

		if _, err := stdin.Write([]byte{0x0D}); err != nil {
			fmt.Println(err)
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	return nil
}
