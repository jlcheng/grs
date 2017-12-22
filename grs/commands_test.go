package grs

import (
	"testing"
	"errors"
	"fmt"
)

// MockRunner holds a sequence of Commands, mapped to their command-line text. When the user specifies a command text,
// it returns the corresponding command and advances to the next command in memory.
type MockRunner struct {
	_commands []*Command
}

func (m *MockRunner) Add(cmd *Command) {
	m._commands = append(m._commands, cmd)
}

func (m *MockRunner) Command(name string, arg ...string) *Command {
	if len(m._commands) == 0 {
		return NewCommandHelper(make([]byte,0), errors.New("no commands configured"))
	}
	r := m._commands[0]
	m._commands = m._commands[1:]
	return r
}

func NewMockRunner() MockRunner {
	m := MockRunner{}
	m._commands = make([]*Command,0)
	return m
}

var echoOne *Command = NewCommandHelper([]byte("one"), nil)
var echoTwo *Command = NewCommandHelper([]byte("two"), nil)
var failed *Command = NewCommandHelper(make([]byte,0), errors.New("failed"))

func TestMockCommandFail(t *testing.T) {
	m := NewMockRunner()
	m.Add(failed)
	cmd := *m.Command("echo","one")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error, got nil")
	}
	if len(out) != 0 {
		t.Errorf("expected empty out, got %v", string(out))
	}
}

func TestMockCommandEmpty(t *testing.T) {
	m := NewMockRunner()
	cmd := *m.Command("echo","one")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error, got nil")
	}
	if fmt.Sprintf("%v", err) != "no commands configured" {
		t.Error("expected error message not found")
	}
	if len(out) != 0 {
		t.Errorf("expected empty out, got %v", string(out))
	}
}

func TestMockCommandOk(t *testing.T) {
	m := NewMockRunner()
	m.Add(echoOne)
	cmd := *m.Command("echo","one")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected ok, got error: %v\n", err)
	}
	if s := string(out); s != "one" {
		t.Errorf("expected 'one', got %v", s)
	}
}

func TestMockCommandMulti(t *testing.T) {
	m := NewMockRunner()
	m.Add(echoOne)
	m.Add(echoTwo)
	m.Add(failed)
	cmd := *m.Command("echo","one")
	out, err := cmd.CombinedOutput()
	if s := string(out); s != "one" {
		t.Errorf("expected 'one', got %v", s)
	}
	cmd = *m.Command("echo","two")
	out, err = cmd.CombinedOutput()
	if s := string(out); s != "two" {
		t.Errorf("expected 'two', got %v", s)
	}
	cmd = *m.Command("invalid")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error, got nil")
	}
	if len(out) != 0 {
		t.Errorf("expected empty out, got %v", string(out))
	}
}
