package process

import (
	"github.com/LonelyPale/goutils/errors"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"strings"
)

var ErrNotFoundProcess = errors.New("process not found")

func FindProcess(i interface{}) (*process.Process, error) {
	switch v := i.(type) {
	case int:
		return process.NewProcess(int32(v))
	case int32:
		return process.NewProcess(v)
	case string:
		procs, err := FindProcessByName(v)
		if err != nil {
			return nil, err
		}
		if len(procs) == 0 {
			return nil, ErrNotFoundProcess
		}
		return procs[0], nil
	default:
		return nil, errors.New("invalid type")
	}
}

func FindProcessByName(name string) ([]*process.Process, error) {
	if len(name) == 0 {
		return nil, errors.New("the name cannot be empty")
	}

	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	res := make([]*process.Process, 0)
	for _, proc := range procs {
		pname, err := proc.Name()
		if err != nil {
			return nil, err
		}

		if pname == name {
			res = append(res, proc)
		}
	}

	return res, nil
}

func FindProcessByCmd(cmd string) ([]*process.Process, error) {
	if len(cmd) == 0 {
		return nil, errors.New("the cmd cannot be empty")
	}

	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	res := make([]*process.Process, 0)
	for _, proc := range procs {
		cmdline, err := proc.Cmdline()
		if err != nil {
			return nil, err
		}

		if strings.Index(cmdline, cmd) > -1 {
			res = append(res, proc)
		}
	}

	return res, nil
}

func CurrentProcess() (*process.Process, error) {
	return process.NewProcess(int32(os.Getpid()))
}
