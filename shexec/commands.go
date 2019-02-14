package shexec

import (
	"os/exec"
)

type Command interface {
	// WithDir sets the working directory for the command.
	WithDir(dir string) Command

	// CombinedOutput runs the command and returns its combined standard output and standard error.
	CombinedOutput() ([]byte, error)
}

type CommandWrapper struct {
	delegate *exec.Cmd
}

func (cmd *CommandWrapper) WithDir(dir string) Command {
	cmd.delegate.Dir = dir
	return cmd
}

func (cmd *CommandWrapper) CombinedOutput() ([]byte, error) {
	return cmd.delegate.CombinedOutput()
}

type CommandRunner interface {
	Command(name string, arg ...string) Command
}

type ExecRunner struct{}

func (r *ExecRunner) Command(name string, arg ...string) Command {
	delegate := exec.Command(name, arg...)
	return &CommandWrapper{delegate}
}
