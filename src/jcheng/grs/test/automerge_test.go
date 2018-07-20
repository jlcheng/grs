package test

import (
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
)

// TestAutoFFMerge_Fail verifies that merge --ff-only is not invoked
func TestAutoFFMerge_Noop(t *testing.T) {
	verifiedGitNotCalled(t, status.DIR_INVALID, status.BRANCH_BEHIND, status.INDEX_UNMODIFIED)

	verifiedGitNotCalled(t, status.DIR_VALID, status.BRANCH_AHEAD, status.INDEX_UNMODIFIED)
	verifiedGitNotCalled(t, status.DIR_VALID, status.BRANCH_UNKNOWN, status.INDEX_UNMODIFIED)

	verifiedGitNotCalled(t, status.DIR_VALID, status.BRANCH_BEHIND, status.INDEX_UNKNOWN)
	verifiedGitNotCalled(t, status.DIR_VALID, status.BRANCH_BEHIND, status.INDEX_MODIFIED)
}

func TestAutoFFMerge_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git merge --ff-only", Ok(""))

	ctx := grs.NewAppContextWithRunner(runner)

	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	repo.Branch = status.BRANCH_BEHIND
	repo.Index = status.INDEX_UNMODIFIED
	script.NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 1 {
		t.Error("git merge not invoked as expected")
	}
}

func verifiedGitNotCalled(t *testing.T, dir status.Dirstat, branch status.Branchstat, index status.Indexstat) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok(""))

	ctx := grs.NewAppContextWithRunner(runner)

	repo := status.NewRepo("")
	repo.Dir = dir
	repo.Branch = branch
	repo.Index = index
	script.NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 0 {
		t.Errorf("unexpected `git merge` when dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}
