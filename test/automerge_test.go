package test

import (
	"testing"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"io/ioutil"
	"path/filepath"
	"time"
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

func TestGetActivityTime(t *testing.T) {
	oldwd, err := os.Getwd()
	d, err := ioutil.TempDir("", "grstest")
	if err != nil {
		t.Fatalf("TestGetActivityTime: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal("TestGetActivityTime.defer: ", err)
		}
		if err := os.RemoveAll(d); err != nil {
			t.Fatal("TestGetActivityTime.defer: ", err)
		}
	}()

	if err := os.Chdir(d); err != nil {
		t.Fatalf("TestGetActivityTime: %v", err)
	}


	os.Mkdir(filepath.Join(d, ".git"), 0777)
	fname := filepath.Join(d, ".git", "HEAD")
	fh, err := os.Create(fname)
	fh.Close()

	atime := time.Date(1900, time.January, 1, 1, 0, 0, 0, time.UTC)
	mtime := time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	if err := os.Chtimes(fname, atime, mtime); err != nil {
		t.Fatalf("TestGetActivityTime: %v", err)
	}

	activity, err := script.GetActivityTime(d)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if !activity.Equal(mtime) {
		t.Error("unexpected last activity time: ", activity)
	}
}

