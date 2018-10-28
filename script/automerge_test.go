package script

import (
	"jcheng/grs/shexec"
	"testing"
)

// TestAutoFFMerge_Fail verifies that merge --ff-only is not invoked
func TestAutoFFMerge_Noop(t *testing.T) {
	verify_AutoFFMerge_NoGitExec(t, DIR_INVALID, BRANCH_BEHIND, INDEX_UNMODIFIED)

	verify_AutoFFMerge_NoGitExec(t, DIR_VALID, BRANCH_AHEAD, INDEX_UNMODIFIED)
	verify_AutoFFMerge_NoGitExec(t, DIR_VALID, BRANCH_UNKNOWN, INDEX_UNMODIFIED)

	verify_AutoFFMerge_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_UNKNOWN)
	verify_AutoFFMerge_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_MODIFIED)
}

func TestAutoFFMerge_Ok(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git merge --ff-only", shexec.Ok(""))

	ctx := shexec.NewAppContextWithRunner(runner)

	repo := NewRepo("")
	repo.Dir = DIR_VALID
	repo.Branch = BRANCH_BEHIND
	repo.Index = INDEX_UNMODIFIED
	NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 1 {
		t.Error("git merge not invoked as expected")
	}
}

func verify_AutoFFMerge_NoGitExec(t *testing.T, dir Dirstat, branch Branchstat, index Indexstat) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git", shexec.Ok(""))

	ctx := shexec.NewAppContextWithRunner(runner)

	repo := NewRepo("")
	repo.Dir = dir
	repo.Branch = branch
	repo.Index = index
	NewScript(ctx, repo).AutoFFMerge()

	if runner.HistoryCount("git merge --ff-only") != 0 {
		t.Errorf("unexpected `git merge` when dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}
