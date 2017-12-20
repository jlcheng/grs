package grs

type Command interface {
	CombinedOutput() ([]byte, error)
}

type CommandRunner interface {
	Command(name string, arg ...string) Command
}

type MockRunner struct {
	commands map[string][]Command
}

func (m *MockRunner) Add(name string, cmd Command) {
	if _, ok := m.commands[name]; !ok {
		m.commands[name] = make([]Command, 0)
	}
	m.commands[name] = append(m.commands[name], cmd)
}

func (m *MockRunner) Command(name string, arg ...string) Command {
	cmd := m.commands[name][0]
	m.commands[name] = m.commands[name][1:]
	return cmd
}

func NewMockRunner() MockRunner {
	m := MockRunner{}
	m.commands = make(map[string][]Command)
	return m
}


type MockCommand struct {
	f func() ([]byte, error)
}

func (m MockCommand) CombinedOutput() ([]byte, error) {
	return m.f()
}

func NewMockCommand(bytes []byte, err error) MockCommand {
	f := func() ([]byte, error) {
		return bytes, err
	}
	return MockCommand{f}
}