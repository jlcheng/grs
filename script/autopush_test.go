package script

import (
	"jcheng/grs/shexec"
	"testing"
)

// Verify that git push does not get called when repo status is unexpected
func TestAutoPush_Noop(t *testing.T) {
	verify_AutoPush_NoGitExec(t, DIR_INVALID, BRANCH_AHEAD, INDEX_UNKNOWN)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_AHEAD, INDEX_MODIFIED)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_UNKNOWN, INDEX_UNMODIFIED)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_UNMODIFIED)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_MODIFIED)
}

// TODO JCHENG
func TestAutoPush_Ok(t *testing.T) {

}

// Given the Dirstat, Branchstat, and Indexstat, signal an error if git push was called
func verify_AutoPush_NoGitExec(t *testing.T, dir Dirstat, branch Branchstat, index Indexstat) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git", shexec.Ok(""))

	ctx := shexec.NewAppContextWithRunner(runner)

	repo := NewRepo("")
	repo.Dir = dir
	repo.Branch = branch
	repo.Index = index
	NewScript(ctx, repo).AutoPush()

	if runner.HistoryCount("git push") != 0 {
		t.Errorf("unexpected `git push` given dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}
