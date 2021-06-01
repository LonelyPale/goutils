package cmd

import (
	"testing"
)

func TestPrint(t *testing.T) {
	Print(Red, "red\n")
}

func TestExecs(t *testing.T) {
	test(t, Execs([]string{
		"pwd",
		"ls /",
		"who",
		"ping -c5 localhost",
	}))
}

func TestExec(t *testing.T) {
	test(t, Exec("ls /"))
	test(t, Exec("pwd"))
	test(t, Exec("ping -c 8 localhost"))
}

func TestShell(t *testing.T) {
	test(t, Shell("ls /"))
}

func TestSudo(t *testing.T) {
	test(t, Sudo("ls /", Options{
		Echo:   true,
		Passwd: "wyb123456",
	}))
}

func TestCommand_ReadAll(t *testing.T) {
	cmd, err := NewCommand("ping -c 8 localhost")
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	str, err := cmd.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(str)
}

func test(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
	println()
}
