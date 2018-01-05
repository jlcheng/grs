package test

import (
	"errors"
	"jcheng/grs/grs"
	"regexp"
	"strings"
	"fmt"
)

type clist struct {
	pattern *regexp.Regexp
	commands []*grs.Command
}

// MockRunner holds a sequence of Commands, mapped to their command-line text. When the user specifies a command text,
// it returns the corresponding command and advances to the next command in memory.
type MockRunner struct {
	_commands []*grs.Command
	commands map[string]*clist
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
	fmt.Printf("%p\n", &v.commands)
	v.commands = append(v.commands, cmd)
	fmt.Println(v.commands)
	m.commands[s] = v
}

func (m *MockRunner) Command(name string, arg ...string) *grs.Command {
	full := strings.Join(append([]string{name}, arg...), " ")
	fmt.Println(m._commands)
	for k := range m.commands {
		if v, ok := m.commands[k]; ok {
			fmt.Println(full, k, ok, v)
			if v.pattern.MatchString(full) && len(v.commands) != 0 {
				r := v.commands[0]
				v.commands = v.commands[1:]
				return r
			}
		}
	}
	if len(m._commands) == 0 {
		return grs.NewCommandHelper(make([]byte,0), errors.New("no commands configured"))
	}
	r := m._commands[0]
	m._commands = m._commands[1:]
	return r
}

func NewMockRunner() *MockRunner {

	return &MockRunner{
		_commands:make([]*grs.Command, 0),
		commands: make(map[string]*clist),
	}
}
