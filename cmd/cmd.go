package cmd

import (
	"bytes"
	"io"
	"sync"
)

const (
	And   = "&&" //与
	Or    = "||" //或
	Alone = ";"  //单独执行
)

func Execs(commands []string, control ...string) error {
	for _, command := range commands {
		if err := Exec(command); err != nil {
			return err
		}
		println()
	}
	return nil
}

func Exec(command string, opts ...Options) error {
	cmd, err := NewCommand(command, opts...)
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if len(opts) > 0 && opts[0].Async {
		return nil
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()
		cmd.Print()
	}()
	go func() {
		defer wg.Done()
		cmd.PrintError()
	}()

	return cmd.Wait()
}

func Shell(command string, opts ...Options) error {
	var buf bytes.Buffer
	buf.WriteString("bash -c '")
	buf.WriteString(command)
	buf.WriteString("'")
	return Exec(buf.String(), opts...)
}

func Sudo(command string, opts ...Options) error {
	var buf bytes.Buffer
	buf.WriteString("sudo -S ")
	buf.WriteString(command)

	cmd, err := NewCommand(buf.String(), opts...)
	if err != nil {
		return err
	}

	if _, err := cmd.In.Write([]byte(opts[0].Passwd + "\r\n")); err != nil {
		Print(Red, err.Error())
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()
		cmd.Print()
	}()
	go func() {
		defer wg.Done()
		for {
			line, err := cmd.ErrorLine()
			if err != nil && err != io.EOF {
				Print(Red, err.Error())
				return
			}

			if line == "Password:" {
				continue
			}

			Print(Red, line)
			if cmd.Exited() && cmd.err.Len() == 0 {
				return
			}
		}
	}()

	return cmd.Wait()
}
