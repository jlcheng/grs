package test

import (
	"testing"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"jcheng/grs/config"
)

// TestGetRepoStatus_Git_Fail verifies that git errors result in BRANCH_UNKNOWN
func TestGetRepoStatus_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetRepoStatus(grs.NewAppContext(), runner, rstat)
	if rstat.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, rstat.Branch)
	}
}

// TestGetRepoStatus_Git_Ok verifies several happy paths
func TestGetRepoStatus_Git_Ok(t *testing.T) {
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

// TestGetRepoStatus_Git_From_Ctx Verifies that the TestGetRepo script gets its 'git' executable from ctx
func TestGetRepoStatus_Git_From_Ctx(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("^/path/to/git", Ok("0\t0\n"))

	ctx := grs.NewAppContext()
	ctx.ConfParams(&config.ConfigParams{User:"data/config.json"})

	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetRepoStatus(ctx, runner, rstat)
	if rstat.Dir == status.DIR_INVALID {
		t.Error("Unexpected rstat.Dir, got DIR_INVALID")
		return
	}
	if rstat.Branch != status.BRANCH_UPTODATE {
		t.Error("Unexpected rstat.Branch, got", rstat.Branch)
		return
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected status.Branchstat) {
	runner := NewMockRunner()
	runner.Add(NewCommandHelper([]byte(out), nil))

	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetRepoStatus(grs.NewAppContext(), runner, rstat)
	got := rstat.Branch
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}