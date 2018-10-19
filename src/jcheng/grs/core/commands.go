package grs

import (
	"os/exec"
)

type Command interface {
	CombinedOutput() ([]byte, error)
}

type CommandRunner interface {
	Command(name string, arg ...string) Command
}

type ExecRunner struct{}

func (r *ExecRunner) Command(name string, arg ...string) Command {
	c := exec.Command(name, arg...)
	return Command(c)
}
