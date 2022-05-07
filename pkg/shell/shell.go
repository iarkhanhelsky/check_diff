package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

type shellOptions struct {
	allowExitCodes []int
}

func (opts shellOptions) checkExitCode(cmd *exec.Cmd) bool {
	for _, c := range opts.allowExitCodes {
		if c == cmd.ProcessState.ExitCode() {
			return true
		}
	}

	return false
}

type ShellOption func(opts *shellOptions)

func AllowExitCodes(codes ...int) ShellOption {
	return func(opts *shellOptions) {
		opts.allowExitCodes = codes
	}
}

type Shell interface {
	Execute(cmd string, args ...string) ([]byte, error)
}

type localShell struct {
	options shellOptions
}

var _ Shell = &localShell{}

func NewLocalShell(opts ...ShellOption) Shell {
	options := shellOptions{}
	for _, o := range opts {
		o(&options)
	}
	return &localShell{options: options}
}

func (localShell *localShell) Execute(command string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil && !localShell.options.checkExitCode(cmd) {
		return nil, fmt.Errorf("failed to execute %s: %v", command, err)
	}

	return stdout.Bytes(), nil
}
