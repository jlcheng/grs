package test

import (
	"testing"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"errors"
)

func TestFetch_Git_Failed(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(NewCommandHelper([]byte(""), errors.New("failed")))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.Fetch(runner, rstat)
	if rstat.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, rstat.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", NewCommandHelper([]byte("0"),nil))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.Fetch(runner, rstat)
	if rstat.Dir == status.DIR_INVALID {
		t.Error("Unexpected rstat.Dir, got DIR_INVALID")
	}
}

