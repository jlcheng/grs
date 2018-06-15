package test

import (
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
)

// TestGetRepoStatus_Git_Fail verifies that git errors result in BRANCH_UNKNOWN
func TestGetRepoStatus_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.GetRepoStatus(grs.NewAppContextWithRunner(runner), repo)
	if repo.Branch != status.BRANCH_UNTRACKED {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNTRACKED, repo.Branch)
	}
}

// TestGetRepoStatus_Git_Ok verifies several happy paths
func TestGetRepoStatus_Git_Ok(t *testing.T) {
	var statustests = []struct {
		out      string
		expected status.Branchstat
	}{
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
	runner.AddMap("^/path/to/git rev-parse", Ok(""))
	runner.AddMap("^/path/to/git rev-list", Ok("0\t0\n"))

	ctx := grs.NewAppContextWithRunner(runner)
	cp := &config.ConfigParams{User: "data/config.json"}
	if conf, _ := config.ReadConfig(cp); conf != nil {
		ctx.SetGitExec(conf.Git)
	}

	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.GetRepoStatus(ctx, repo)
	if repo.Dir == status.DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
		return
	}
	if repo.Branch != status.BRANCH_UPTODATE {
		t.Error("Unexpected repo.Branch, got", repo.Branch)
		return
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected status.Branchstat) {
	runner := NewMockRunner()
	runner.AddMap("git rev-parse", Ok("..."))
	runner.AddMap("git rev-list", Ok(out))

	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.GetRepoStatus(grs.NewAppContextWithRunner(runner), repo)
	got := repo.Branch
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}
