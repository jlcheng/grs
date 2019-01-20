
package script

import (
	"jcheng/grs/shexec"
	"testing"
)

// TestGetRepoStatus_Git_Fail verifies that git errors result in BRANCH_UNKNOWN
func TestGetRepoStatus_Git_Fail(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.Add(shexec.Error("failed"))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(shexec.NewAppContextWithRunner(runner), repo).GetRepoStatus()
	if repo.Branch != BRANCH_UNTRACKED {
		t.Errorf("expected %s, got: %v\n", BRANCH_UNTRACKED, repo.Branch)
	}
}

// TestGetRepoStatus_Git_Ok verifies several happy paths
func TestGetRepoStatus_Git_Ok(t *testing.T) {
	var statustests = []struct {
		out      string
		expected Branchstat
	}{
		{"0\t1\n", BRANCH_AHEAD},
		{"1\t0\n", BRANCH_BEHIND},
		{"1\t1\n", BRANCH_DIVERGED},
		{"invalid\n", BRANCH_UNKNOWN},
	}
	for _, elem := range statustests {
		helpGetRepoStatus(t, elem.out, elem.expected)
	}
}

// TestGetRepoStatus_Git_From_Ctx Verifies that the TestGetRepo script gets its 'git' executable from ctx
func TestGetRepoStatus_Git_From_Ctx(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("^/path/to/git rev-parse", shexec.Ok(""))
	runner.AddMap("^/path/to/git rev-list", shexec.Ok("0\t0\n"))

	ctx := shexec.NewAppContextWithRunner(runner)
	ctx.SetGitExec("/path/to/git")

	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(ctx, repo).GetRepoStatus()
	if repo.Dir == DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
		return
	}
	if repo.Branch != BRANCH_UPTODATE {
		t.Error("Unexpected repo.Branch, got", repo.Branch)
		return
	}
}

func TestGetCommitTime(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("^git log -1 --format=%cr", shexec.Ok("5 minutes ago\n"))
	ctx := shexec.NewAppContextWithRunner(runner)

	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(ctx, repo).GetCommitTime()
	expected := "5 minutes ago"
	if repo.CommitTime != expected {
		t.Error("Unexpected commit time, got [" + repo.CommitTime + "]")
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected Branchstat) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git rev-parse", shexec.Ok("..."))
	runner.AddMap("git rev-list", shexec.Ok(out))

	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(shexec.NewAppContextWithRunner(runner), repo).GetRepoStatus()
	got := repo.Branch
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}
