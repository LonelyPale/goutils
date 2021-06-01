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
	}))
}

func TestExec(t *testing.T) {
	test(t, Exec("ls /"))
	test(t, Exec("pwd"))
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

func test(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
