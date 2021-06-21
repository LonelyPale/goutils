package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
	"golang.org/x/crypto/ssh"

	"github.com/LonelyPale/goutils/errors"
)

var (
	Red     = color.New(color.FgRed)
	Blue    = color.New(color.FgBlue)
	Green   = color.New(color.FgGreen)
	Yellow  = color.New(color.FgYellow)
	Magenta = color.New(color.FgMagenta)
	Cyan    = color.New(color.FgCyan)
	White   = color.New(color.FgWhite)
	Black   = color.New(color.FgBlack)
)

type Buffer struct {
	buf bytes.Buffer
	mu  sync.RWMutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.Read(p)
}

func (b *Buffer) ReadString(delim byte) (line string, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.ReadString(delim)
}

func (b *Buffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.String()
}

func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.Len()
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *Buffer) Copy(src io.Reader) (written int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return io.Copy(&b.buf, src)
}

type Options struct {
	Stdin  io.ReadWriter
	Stdout io.ReadWriter
	Stderr io.ReadWriter
	Echo   bool
	Async  bool
	Passwd string
}

type Command struct {
	Cmd  *exec.Cmd
	Ssh  *ssh.Session
	In   io.Writer
	Out  io.Reader
	Err  io.Reader
	out  Buffer
	err  Buffer
	opts Options
}

func NewCommand(command string, opts ...Options) (*Command, error) {
	wrapper := new(Command)

	splitArgs, err := shellquote.Split(command)
	if err != nil {
		return nil, err
	}
	numArguments := len(splitArgs) - 1
	if numArguments < 0 {
		return nil, errors.New("cmd: No command given to new")
	}
	path, err := exec.LookPath(splitArgs[0])
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	if numArguments >= 1 {
		cmd = exec.Command(path, splitArgs[1:]...)
	} else {
		cmd = exec.Command(path)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdin, stdout, stderr io.ReadWriter
	if len(opts) > 0 {
		if opts[0].Stdin == nil {
			opts[0].Stdin = &Buffer{}
		}
		if opts[0].Stdout == nil {
			opts[0].Stdout = &Buffer{}
		}
		if opts[0].Stderr == nil {
			opts[0].Stderr = &Buffer{}
		}

		stdin = opts[0].Stdin
		stdout = opts[0].Stdout
		stderr = opts[0].Stderr
		wrapper.opts = opts[0]
	} else {
		stdin = &Buffer{}
		stdout = &Buffer{}
		stderr = &Buffer{}
		wrapper.opts = Options{
			Echo: true,
		}
	}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	wrapper.In = stdin
	wrapper.Out = stdout
	wrapper.Err = stderr
	wrapper.Cmd = cmd
	return wrapper, nil
}

func (c *Command) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *Command) Start() error {
	if c.opts.Echo {
		Print(Blue, c.Cmd.String()+"\r\n")
	}
	return c.Cmd.Start()
}

func (c *Command) Wait() error {
	return c.Cmd.Wait()
}

func (c *Command) Close() error {
	return c.Cmd.Process.Kill()
}

func (c *Command) Kill() error {
	pid := c.Cmd.Process.Pid
	if err := syscall.Kill(-1*pid, syscall.SIGKILL); err != nil {
		return err
	}
	return nil
}

func (c *Command) Exited() bool {
	return c.Cmd.ProcessState != nil && c.Cmd.ProcessState.Exited()
}

func (c *Command) Write(data []byte) (int, error) {
	return c.In.Write(data)
}

func (c *Command) WriteLine(command string) error {
	_, err := io.WriteString(c.In, command+"\r\n")
	return err
}

func (c *Command) Read(data []byte) (int, error) {
	return c.Out.Read(data)
}

func (c *Command) ReadLine() (string, error) {
	if _, err := c.out.Copy(c.Out); err != nil { //不返回io.EOF
		return "", err
	}
	return c.out.ReadString('\n') //会返回io.EOF
}

func (c *Command) ReadAll() (string, error) {
	for {
		if _, err := c.out.Copy(c.Out); err != nil {
			return c.out.String(), err
		}

		if c.Exited() {
			return c.out.String(), nil
		}
	}
}

func (c *Command) Error(data []byte) (int, error) {
	return c.Err.Read(data)
}

func (c *Command) ErrorLine() (string, error) {
	if _, err := c.err.Copy(c.Err); err != nil {
		return "", err
	}
	return c.err.ReadString('\n')
}

func (c *Command) ErrorAll() (string, error) {
	for {
		if _, err := c.err.Copy(c.Err); err != nil {
			return c.err.String(), err
		}

		if c.Exited() {
			return c.err.String(), nil
		}
	}
}

func (c *Command) Output() (string, error) {
	if err := c.Run(); err != nil {
		return "", err
	}
	return c.ReadAll()
}

func (c *Command) CombinedOutput() (string, error) {
	if err := c.Run(); err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	s := bufio.NewScanner(io.MultiReader(c.Out, c.Err))
	for s.Scan() {
		buffer.Write(s.Bytes())
	}
	return buffer.String(), s.Err()
}

func (c *Command) Print() {
	for {
		line, err := c.ReadLine()
		if err != nil && err != io.EOF {
			Print(Red, err.Error())
			return
		}

		Print(Green, line)
		if c.Exited() && c.out.Len() == 0 {
			return
		}
	}
}

func (c *Command) PrintError() {
	for {
		line, err := c.ErrorLine()
		if err != nil && err != io.EOF {
			Print(Red, err.Error())
			return
		}

		Print(Red, line)
		if c.Exited() && c.err.Len() == 0 {
			return
		}
	}
}

func Print(c *color.Color, s string) {
	if len(s) == 0 {
		return
	}
	if _, err := c.Print(s); err != nil {
		fmt.Println(err)
	}
}
