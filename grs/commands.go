package grs

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
