package shexec

import (
	"testing"
)


// Verifies that the default GetGitExec() is `git`
func TestGetGitExecDefault(t *testing.T) {
	ctx := NewAppContext()

	if r := ctx.GitExec; r != "git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}
