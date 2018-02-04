package grs

import (
	"os/exec"
	"jcheng/grs/compat"
)

type Command interface {
	CombinedOutput() ([]byte, error)
}

type CommandRunner interface {
	Command(name string, arg ...string) *Command
}

type ExecRunner struct { }

func (r ExecRunner) Command(name string, arg ...string) *Command {
	c := exec.Command(name, arg...)
	compat.BeforeCmd(c)
	ref := Command(c)
	return &ref
}
