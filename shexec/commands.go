package shexec

import (
	"fmt"
	"os"
	"os/exec"
)

type Command interface {
	// WithDir sets the working directory for the command.
	WithDir(dir string) Command

	// CombinedOutput runs the command and returns its combined standard output and standard error.
	CombinedOutput() ([]byte, error)
}

// CommandWrapper implements the shexec.Command interface
type CommandWrapper struct {
	delegate *exec.Cmd
}

// WithDir sets the working directory for the command.
func (cmd *CommandWrapper) WithDir(dir string) Command {
	cmd.delegate.Dir = dir
	return cmd
}

// CombinedOutput runs the command and returns its combined standard output and standard error.
func (cmd *CommandWrapper) CombinedOutput() ([]byte, error) {
	b, err := cmd.delegate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(string(b))
	}
	return b, nil
}

// CommandRunner resolves an implmentation of Command
type CommandRunner interface {
	// Command creates an instance of a Command object
	Command(name string, arg ...string) Command
}

// ExecRunner is an implementation of CommandRunner that delegates to os/exec.Command
type ExecRunner struct{}

// Command creates an instance of a Command object
func (r *ExecRunner) Command(name string, arg ...string) Command {
	delegate := exec.Command(name, arg...)
	envcopy := make([]string, len(os.Environ()))
	copy(envcopy, os.Environ())
	delegate.Env = append(envcopy, "GIT_SSH_COMMAND=ssh -o BatchMode=yes")
	return &CommandWrapper{delegate}
}
