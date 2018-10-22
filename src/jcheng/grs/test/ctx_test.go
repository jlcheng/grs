package test

import (
	"jcheng/grs/core"
	"testing"
)


// Verifies that the default GetGitExec() is `git`
func TestGetGitExecDefault(t *testing.T) {
	ctx := grs.NewAppContext()

	if r := ctx.GetGitExec(); r != "git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}
