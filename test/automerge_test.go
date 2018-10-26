package test

import (
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"testing"
)

// TestAutoFFMerge_Fail verifies that merge --ff-only is not invoked
func TestAutoFFMerge_Noop(t *testing.T) {
	verifiedGitNotCalled(t, script.DIR_INVALID, script.BRANCH_BEHIND, script.INDEX_UNMODIFIED)

	verifiedGitNotCalled(t, script.DIR_VALID, script.BRANCH_AHEAD, script.INDEX_UNMODIFIED)
	verifiedGitNotCalled(t, script.DIR_VALID, script.BRANCH_UNKNOWN, script.INDEX_UNMODIFIED)

	verifiedGitNotCalled(t, script.DIR_VALID, script.BRANCH_BEHIND, script.INDEX_UNKNOWN)
	verifiedGitNotCalled(t, script.DIR_VALID, script.BRANCH_BEHIND, script.INDEX_MODIFIED)
}

func TestAutoFFMerge_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git merge --ff-only", Ok(""))

	ctx := shexec.NewAppContextWithRunner(runner)

	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	repo.Branch = script.BRANCH_BEHIND
	repo.Index = script.INDEX_UNMODIFIED
	script.NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 1 {
		t.Error("git merge not invoked as expected")
	}
}

func verifiedGitNotCalled(t *testing.T, dir script.Dirstat, branch script.Branchstat, index script.Indexstat) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok(""))

	ctx := shexec.NewAppContextWithRunner(runner)

	repo := script.NewRepo("")
	repo.Dir = dir
	repo.Branch = branch
	repo.Index = index
	script.NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 0 {
		t.Errorf("unexpected `git merge` when dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}
