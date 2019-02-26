package script

import (
	"jcheng/grs/shexec"
	"strings"
	"testing"
	"time"
)

func TestAutoPushGenCommitMsg(t *testing.T) {
	nowRetval, err := time.Parse(time.RFC3339, "1234-05-06T07:08:09Z")
	if err != nil {
		t.Error(err)
	}
	clock := &MockClock{NowRetval: nowRetval}
	if got := AutoPushGenCommitMsg(clock); !strings.Contains(got, "1234-05-06T07:08:09Z") {
		t.Error("expected timestamp missing. got:", got)
	}

}

// Verify that git push does not get called when repo status is unexpected
func TestAutoPush_Noop(t *testing.T) {
	verify_AutoPush_NoGitExec(t, DIR_INVALID, BRANCH_AHEAD, INDEX_UNKNOWN)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_AHEAD, INDEX_MODIFIED)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_UNKNOWN, INDEX_UNMODIFIED)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_UNMODIFIED)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_MODIFIED)
}

func TestAutoPush_Ok(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git add -A", shexec.Ok(""))
	runner.AddMap("git commit -m ", shexec.Ok(""))
	runner.AddMap("git push", shexec.Ok(""))

	ctx := NewAppContext(WithCommandRunner(runner))

	repo := NewRepo("")
	repo.Dir = DIR_VALID
	repo.Branch = BRANCH_AHEAD
	repo.Index = INDEX_MODIFIED
	repo.PushAllowed = true
	NewScript(ctx, repo).AutoPush()

	if runner.HistoryCount("git commit -m") != 1 {
		t.Error("git commit not invoked as expected")
	}
	if runner.HistoryCount("git push") != 1 {
		t.Error("git push not invoked as expected")
	}
}

// Given the Dirstat, Branchstat, and Indexstat, signal an error if git push was called
func verify_AutoPush_NoGitExec(t *testing.T, dir Dirstat, branch Branchstat, index Indexstat) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git", shexec.Ok(""))

	ctx := NewAppContext(WithCommandRunner(runner))

	repo := NewRepo("")
	repo.Dir = dir
	repo.Branch = branch
	repo.Index = index
	NewScript(ctx, repo).AutoPush()

	if runner.HistoryCount("git push") != 0 {
		t.Errorf("unexpected `git push` given dirstat=%v, branchstat=%v, indexstat=%v\n", dir, branch, index)
	}
}

/*
commits *modified and pushes it upstream

After setup, before run
A--B            source
    \
     *modified  dest

After AutoPush(), *modified turned into commit C and is pushed to source

A--B--C        source
       \
        C      dest
*/
func TestAutoPush_IT_Test_1(t *testing.T) {
	const TEST_ID = "TestAutoPush_IT_Test_1"
	tmpdir, cleanup := MkTmpDir1(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	repo := NewRepo(gh.Getwd())
	repo.PushAllowed = true
	s := NewScript(NewAppContext(), repo)
	s.BeforeScript()

	gh.TouchAndCommit("A.txt", "Commit A")
	gh.TouchAndCommit("B.txt", "Commit B")
	gh.Exec(gh.Git(), "push")
	gh.Touch("C.txt")

	s.Update()

	if repo.Branch != BRANCH_UPTODATE || repo.Index != INDEX_MODIFIED {
		t.Fatal(TEST_ID + ": setup failed")
	}

	s.AutoPush()

	if repo.Branch != BRANCH_UPTODATE || repo.Index != INDEX_UNMODIFIED {
		t.Fatal(TEST_ID + ": got repo was not UPTODATE and UNMODIEFIED")
	}
}
