package process

import (
	"fmt"
	"os"
	"testing"

	"github.com/shirou/gopsutil/v3/process"
)

func TestPs(t *testing.T) {
	name, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name:", name)

	procs, err := process.Processes()
	if err != nil {
		t.Fatal(err)
	}

	for _, proc := range procs {
		pname, err := proc.Name()
		if err != nil {
			t.Fatal(err)
		}
		if pname == "ping" {
			t.Log(pname)
			t.Log(proc.Cmdline())
			t.Log(proc.Status())
			t.Log(proc.Exe())
			t.Log(proc.Background())
			t.Log(proc.Foreground())
			t.Log(proc.CmdlineSlice())
			t.Log(proc.Cwd())
		}
	}
}

func TestFindProcess(t *testing.T) {
	var pid int = 10
	proc, err := FindProcess(pid)
	if err != nil {
		if err == ErrorProcessNotRunning {
			t.Error(err)
		} else {
			t.Fatal(err)
		}
	} else {
		printProcess(proc)
	}

	pname := "timed"
	proc, err = FindProcess(pname)
	if err != nil {
		t.Fatal(err)
	}
	printProcess(proc)
}

func TestFindProcessByCmd(t *testing.T) {
	procs, err := FindProcessByCmd("localhost")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("len: ", len(procs))
	for _, proc := range procs {
		t.Log(proc.Cmdline())
	}
}

func TestCurrentProcess(t *testing.T) {
	proc, err := CurrentProcess()
	if err != nil {
		t.Fatal(err)
	}
	printProcess(proc)
}

func printProcess(proc *process.Process) {
	name, err := proc.Name()
	if err != nil {
		panic(err)
	}
	fmt.Println(proc.Pid, name)
}
