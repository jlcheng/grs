package script

import (
	"jcheng/grs/shexec"
	"testing"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.Add(shexec.Error("failed"))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	ctx := shexec.NewAppContext(shexec.WithCommandRunner(runner))
	NewScript(ctx, repo).Fetch()
	if repo.Branch != BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", BRANCH_UNKNOWN, repo.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git", shexec.Ok("0"))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	ctx := shexec.NewAppContext(shexec.WithCommandRunner(runner))
	NewScript(ctx, repo).Fetch()
	if repo.Dir == DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
	}
}
