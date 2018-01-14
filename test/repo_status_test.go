package test

import (
	"testing"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"os"
	"jcheng/grs/status"
	"errors"
)

func TestGetRepoStatus_Git_Failed(t *testing.T) {
	runner := NewMockRunner()
	var repo grs.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:cwd}
	}
	runner.Add(grs.NewCommandHelper([]byte(""), errors.New("failed")))
	s := script.GetRepoStatus(repo, runner)
	if s.Branch != status.BRANCH_UNKNOWN {
		t.Error("expected %s, got: %v", status.BRANCH_UNKNOWN, s.Branch)
	}
}

func TestGetRepoStatus_Git(t *testing.T) {
	var statustests = []struct {
		out string
		expected status.Branchstat
	} {
		{"0\t1\n", status.BRANCH_AHEAD},
		{"1\t0\n", status.BRANCH_BEHIND},
		{"1\t1\n", status.BRANCH_DIVERGED},
		{"invalid\n", status.BRANCH_UNKNOWN},
	}
	for _, elem := range statustests {
		helpGetRepoStatus(t, elem.out, elem.expected)
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected status.Branchstat) {
	runner := NewMockRunner()
	runner.Add(grs.NewCommandHelper([]byte(out), nil))
	var repo grs.Repo
	if d, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:d}
	}
	got := script.GetRepoStatus(repo, runner).Branch
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}