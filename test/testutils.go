package test

import (
	"errors"
	"jcheng/grs/grs"
	"regexp"
	"strings"
)

type clist struct {
	pattern *regexp.Regexp
	commands []*grs.Command
}

// CommandHelper allows one to easily mock a CommandRunner interface. Useful when you wan to return the same output
// given any input.
type CommandHelper struct {
	f func() ([]byte, error)
}

func (m CommandHelper) CombinedOutput() ([]byte, error) {
	return m.f()
}

func NewCommandHelper(bytes []byte, err error) *grs.Command {
	f := func() ([]byte, error) {
		return bytes, err
	}
	var r grs.Command = CommandHelper{f}
	return &r
}


var _EMPTY_BYTES = []byte("")
// Error is a convenience function for mocking common errors
func Error(msg string) *grs.Command {
	return NewCommandHelper(_EMPTY_BYTES, errors.New(msg))
}

func Ok(msg string) *grs.Command {
	return NewCommandHelper([]byte(msg), nil)
}


// MockRunner holds a sequence of Commands, mapped to their command-line text. When the user specifies a command text,
// it returns the corresponding command and advances to the next command in memory.
type MockRunner struct {
	_commands []*grs.Command
	commands map[string]*clist
	history []string
}

func (m *MockRunner) Add(cmd *grs.Command) {
	m._commands = append(m._commands, cmd)
}

func (m *MockRunner) AddMap(s string, cmd *grs.Command) {
	v, ok := m.commands[s]
	if !ok {
		v = &clist{
			pattern: regexp.MustCompile(s),
			commands: make([]*grs.Command, 0),
		}
	}
	v.commands = append(v.commands, cmd)
	m.commands[s] = v
}

func (m *MockRunner) Command(name string, arg ...string) *grs.Command {
	full := strings.Join(append([]string{name}, arg...), " ")
	m.history = append(m.history, full)

	if len(m._commands) == 0 && len(m.commands) == 0 {
		return NewCommandHelper(make([]byte,0), errors.New("no commands configured"))
	}

	for k := range m.commands {
		if v, ok := m.commands[k]; ok {
			if v.pattern.MatchString(full) && len(v.commands) != 0 {
				r := v.commands[0]
				v.commands = v.commands[1:]
				return r
			}
		}
	}
	if len(m._commands) == 0 {
		return NewCommandHelper(make([]byte,0), errors.New("mock has no commands that match: " + name))
	}

	r := m._commands[0]
	m._commands = m._commands[1:]
	return r
}

// Get a count of matching commands from the runner
func (m *MockRunner) HistoryCount(command string) int {
	var ret int = 0
	p := regexp.MustCompile(command)
	for _, elem := range m.history {
		if p.MatchString(elem) {
			ret++
		}
	}
	return ret
}

func NewMockRunner() *MockRunner {

	return &MockRunner{
		_commands:make([]*grs.Command, 0),
		commands: make(map[string]*clist),
		history:make([]string, 0),
	}
}
