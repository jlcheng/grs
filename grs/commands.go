package grs

import "os/exec"

type Command interface {
	CombinedOutput() ([]byte, error)
}

type CommandRunner interface {
	Command(name string, arg ...string) *Command
}

type CommandHelper struct {
	f func() ([]byte, error)
}

func (m CommandHelper) CombinedOutput() ([]byte, error) {
	return m.f()
}

func NewCommandHelper(bytes []byte, err error) *Command {
	f := func() ([]byte, error) {
		return bytes, err
	}
	var r Command = CommandHelper{f}
	return &r
}

type ExecRunner struct { }

func (r ExecRunner) Command(name string, arg ...string) *Command {
	var c Command = exec.Command(name, arg...)
	return &c
}
