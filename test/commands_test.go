package test

import (
	"testing"
	"errors"
	"fmt"
	"jcheng/grs/grs"
)

var echoOne *grs.Command = grs.NewCommandHelper([]byte("one"), nil)
var echoTwo *grs.Command = grs.NewCommandHelper([]byte("two"), nil)
var failed *grs.Command = grs.NewCommandHelper(make([]byte,0), errors.New("failed"))

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