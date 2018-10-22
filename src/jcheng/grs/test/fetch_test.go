package test

import (
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, repo.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Dir == status.DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
	}
}
