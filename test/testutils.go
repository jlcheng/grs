package test

import (
	"errors"
	"jcheng/grs/grs"
)

// MockRunner holds a sequence of Commands, mapped to their command-line text. When the user specifies a command text,
// it returns the corresponding command and advances to the next command in memory.
type MockRunner struct {
	_commands []*grs.Command
}

func (m *MockRunner) Add(cmd *grs.Command) {
	m._commands = append(m._commands, cmd)
}

func (m *MockRunner) Command(name string, arg ...string) *grs.Command {
	if len(m._commands) == 0 {
		return grs.NewCommandHelper(make([]byte,0), errors.New("no commands configured"))
	}
	r := m._commands[0]
	m._commands = m._commands[1:]
	return r
}

func NewMockRunner() *MockRunner {
	return &MockRunner{_commands:make([]*grs.Command, 0)}
}
