package cmd

import "sync"

func Exec(command string, opts ...Options) error {
	cmd, err := NewCommand(command, opts...)
	if err != nil {
		return err
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
		cmd.PrintError()
	}()

	return cmd.Wait()
}

func Shell(command string, opts ...Options) error {
	return Exec("bash -c"+command, opts...)
}

func Sudo(command string, opts ...Options) error {

	return Exec("sudo", opts...)
}
