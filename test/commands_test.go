package test

import (
	"testing"
	"errors"
	"fmt"
)

var echoOne = NewCommandHelper([]byte("one"), nil)
var echoTwo = NewCommandHelper([]byte("two"), nil)
var dateS = NewCommandHelper([]byte("1515196992"), nil)
var failed = NewCommandHelper(make([]byte,0), errors.New("failed"))

func TestMockCommand_Fail(t *testing.T) {
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

func TestMockCommand_Empty(t *testing.T) {
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

func TestMockCommand_Ok(t *testing.T) {
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

func TestMockCommand_Multi_Ok(t *testing.T) {
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

func TestMockCommandMap_Ok(t *testing.T) {
	m := NewMockRunner()
	m.AddMap("date \\+%s", dateS) // Must to escape + as arg is a regexp
	m.AddMap("echo one", echoOne)

	cmd := *m.Command("echo", "one")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected ok, got error: %v\n", err)
	}
	if s := string(out); s != "one" {
		t.Errorf("expected 'one', got %v", s)
	}

	cmd = *m.Command("date", "+%s")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected ok, got error: %v\n", err)
	}
	if s := string(out); s != "1515196992" {
		t.Errorf("expected '1515196992', got %v", s)
	}
}

func TestMockCommand_HistoryCount(t *testing.T) {
	m := *NewMockRunner()

	m.Command("foo")
	m.Command("fab", "foz")
	m.Command("foo", "fab", "foz")

	var c int
	if c = m.HistoryCount("^foo$"); c != 1 {
		t.Error("expected count(^foo$) == 1, got", c)
	}
	if c = m.HistoryCount("^fab foz$"); c != 1 {
		t.Error("expected count(^fab foz$) == 1, got", c)
	}
	if c = m.HistoryCount("foo"); c != 2 {
		t.Error("expected count(foo) == 2, got", c)
	}
}