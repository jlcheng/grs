package test

import (
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"testing"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	ctx := shexec.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Branch != script.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", script.BRANCH_UNKNOWN, repo.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	ctx := shexec.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Dir == script.DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
	}
}
