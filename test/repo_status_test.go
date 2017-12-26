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
	if d, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:d}
	}
	runner.Add(grs.NewCommandHelper([]byte(""), errors.New("failed")))
	s := script.GetRepoStatus(repo, runner)
	if s != status.UNKNOWN {
		t.Error("expected UNKNOWN, got: %v", s)
	}
}

func TestGetRepoStatus_Git(t *testing.T) {
	var statustests = []struct {
		out string
		expected status.RepoStatus
	} {
		{"0\t1\n", status.AHEAD},
		{"1\t0\n", status.BEHIND},
		{"1\t1\n", status.DIVERGED},
		{"invalid\n", status.UNKNOWN},
	}
	for _, elem := range statustests {
		helpGetRepoStatus(t, elem.out, elem.expected)
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected status.RepoStatus) {
	runner := NewMockRunner()
	runner.Add(grs.NewCommandHelper([]byte(out), nil))
	var repo grs.Repo
	if d, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:d}
	}
	got := script.GetRepoStatus(repo, runner)
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}