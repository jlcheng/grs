package test

import (
	"testing"
	"jcheng/grs/script"
	"jcheng/grs/status"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.Fetch(runner, rstat)
	if rstat.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, rstat.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.Fetch(runner, rstat)
	if rstat.Dir == status.DIR_INVALID {
		t.Error("Unexpected rstat.Dir, got DIR_INVALID")
	}
}

