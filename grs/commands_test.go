package grs

import (
	"testing"
	"errors"
)

func TestMockCommand(t *testing.T) {
	m := NewMockRunner()
	m.Add("echo", NewMockCommand(make([]byte,0), errors.New("failed")))
	cmd := m.Command("echo","1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error, got nil")
	}
	if len(out) != 0 {
		t.Errorf("expected empty out, got %v", string(out))
	}

}
