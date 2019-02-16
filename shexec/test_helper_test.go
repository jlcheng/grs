package shexec

import "testing"

func TestMockRunner(t *testing.T) {
	var _ CommandRunner = &MockRunner{}
}
