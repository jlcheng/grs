package shexec

import (
	"errors"
	"regexp"
	"strings"
)

// CommandHelper allows one to easily mock a CommandRunner interface. Useful when you wan to return the same output
// given any input.
type CommandHelper struct {
	f func() ([]byte, error)
}

func (m *CommandHelper) WithDir(dir string) Command {
	// no-op
	return m
}

func (m *CommandHelper) CombinedOutput() ([]byte, error) {
	return m.f()
}

func NewCommandHelper(bytes []byte, err error) Command {
	f := func() ([]byte, error) {
		return bytes, err
	}
	return &CommandHelper{f}
}

var _EMPTY_BYTES = []byte("")

// Error is a convenience function for mocking common errors
func Error(msg string) Command {
	return NewCommandHelper(_EMPTY_BYTES, errors.New(msg))
}

func Ok(msg string) Command {
	return NewCommandHelper([]byte(msg), nil)
}

// MockRunner holds a sequence of Commands, mapped to their command-line text. When the user specifies a command text,
// it returns the corresponding Command and advances to the next Command.
type MockRunner struct {
	_commands []Command
	commands  map[string]*clist
	history   []string
}

func (m *MockRunner) Add(cmd Command) {
	m._commands = append(m._commands, cmd)
}

func (m *MockRunner) AddMap(s string, cmd Command) {
	v, ok := m.commands[s]
	if !ok {
		v = &clist{
			pattern:  regexp.MustCompile(s),
			commands: make([]Command, 0),
		}
	}
	v.commands = append(v.commands, cmd)
	m.commands[s] = v
}

func (m *MockRunner) Command(name string, arg ...string) Command {
	full := strings.Join(append([]string{name}, arg...), " ")
	m.history = append(m.history, full)

	if len(m._commands) == 0 && len(m.commands) == 0 {
		return Error("no commands configured")
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
		return Error("mock has no commands that match: "+ name)
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
		_commands: make([]Command, 0),
		commands:  make(map[string]*clist),
		history:   make([]string, 0),
	}
}

type clist struct {
	pattern  *regexp.Regexp
	commands []Command
}