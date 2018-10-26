
package test

import (
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"testing"
)

// TestGetRepoStatus_Git_Fail verifies that git errors result in BRANCH_UNKNOWN
func TestGetRepoStatus_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetRepoStatus()
	if repo.Branch != script.BRANCH_UNTRACKED {
		t.Errorf("expected %s, got: %v\n", script.BRANCH_UNTRACKED, repo.Branch)
	}
}

// TestGetRepoStatus_Git_Ok verifies several happy paths
func TestGetRepoStatus_Git_Ok(t *testing.T) {
	var statustests = []struct {
		out      string
		expected script.Branchstat
	}{
		{"0\t1\n", script.BRANCH_AHEAD},
		{"1\t0\n", script.BRANCH_BEHIND},
		{"1\t1\n", script.BRANCH_DIVERGED},
		{"invalid\n", script.BRANCH_UNKNOWN},
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

	ctx := shexec.NewAppContextWithRunner(runner)
	ctx.SetGitExec("/path/to/git")

	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(ctx, repo).GetRepoStatus()
	if repo.Dir == script.DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
		return
	}
	if repo.Branch != script.BRANCH_UPTODATE {
		t.Error("Unexpected repo.Branch, got", repo.Branch)
		return
	}
}

func helpGetRepoStatus(t *testing.T, out string, expected script.Branchstat) {
	runner := NewMockRunner()
	runner.AddMap("git rev-parse", Ok("..."))
	runner.AddMap("git rev-list", Ok(out))

	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetRepoStatus()
	got := repo.Branch
	if got != expected {
		t.Errorf("expected [%v], got [%v]\n", expected, got)
	}
}
