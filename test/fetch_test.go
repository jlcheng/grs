package test

import (
	"testing"
	"jcheng/grs/grs"
	"os"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"errors"
)

func TestFetch_Git_Failed(t *testing.T) {
	runner := NewMockRunner()
	var repo grs.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:cwd}
	}
	runner.Add(NewCommandHelper([]byte(""), errors.New("failed")))
	s := script.Fetch(repo, runner)
	if s.Branch != status.BRANCH_UNKNOWN {
		t.Error("expected %s, got: %v", status.BRANCH_UNKNOWN, s.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", NewCommandHelper([]byte("0"),nil))

	var repo grs.Repo
	if d, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:d}
	}
	rstat := script.Fetch(repo, runner)
	if rstat.Dir == status.DIR_INVALID {
		t.Error("Unexpected rstat.Dir, got DIR_INVALID")
		return
	}
}

