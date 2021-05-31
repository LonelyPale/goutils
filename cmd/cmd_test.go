package cmd

import (
	"testing"
)

func TestExec(t *testing.T) {
	//test(t, Exec("ls /"))
	test(t, Exec("pwd"))
}

func test(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
