package test

import (
	"testing"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
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

	ctx := grs.NewAppContext()

	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	rstat.Branch = status.BRANCH_BEHIND
	rstat.Index = status.INDEX_UNMODIFIED
	script.AutoFFMerge(ctx, runner, rstat)

	if runner.HistoryCount("git merge --ff-only") != 1 {
		t.Error("git merge not invoked as expected")
	}
}

func verifiedGitNotCalled(t *testing.T, dir status.Dirstat, branch status.Branchstat, index status.Indexstat) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok(""))

	ctx := grs.NewAppContext()

	rstat := status.NewRStat()
	rstat.Dir = dir
	rstat.Branch = branch
	rstat.Index = index
	script.AutoFFMerge(ctx, runner, rstat)

	if runner.HistoryCount("git merge --ff-only") != 0 {
		t.Errorf("unexpected `git merge` when dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}