package script

import (
	"jcheng/grs/shexec"
	"testing"
)

func TestGetIndexStatus_Ls_Files_Fail(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git ls-files", shexec.Error("failed"))
	runner.AddMap("git diff-index", shexec.Ok(""))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(NewAppContext(WithCommandRunner(runner)), repo).GetIndexStatus()
	if repo.Index != INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Diff_Index_Fail(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git ls-files", shexec.Ok(""))
	runner.AddMap("git diff-index", shexec.Error("failed"))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(NewAppContext(WithCommandRunner(runner)), repo).GetIndexStatus()
	if repo.Index != INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Unmodified_Ok(t *testing.T) {
	runner := shexec.NewMockRunner()
	runner.AddMap("git ls-files", shexec.Ok(""))
	runner.AddMap("git diff-index", shexec.Ok(""))
	repo := NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(NewAppContext(WithCommandRunner(runner)), repo).GetIndexStatus()
	if repo.Index != INDEX_UNMODIFIED {
		t.Errorf("expected %s, got: %v\n", INDEX_UNMODIFIED, repo.Index)
	}
}

func TestGetIndexStatus_Modified_Ok(t *testing.T) {
	var runner *shexec.MockRunner
	var repo *Repo
	runner = shexec.NewMockRunner()
	runner.AddMap("git ls-files", shexec.Ok("foo\n"))
	runner.AddMap("git diff-index", shexec.Ok(""))
	repo = NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(NewAppContext(WithCommandRunner(runner)), repo).GetIndexStatus()
	if repo.Index != INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", INDEX_MODIFIED, repo.Index)
	}

	runner = shexec.NewMockRunner()
	runner.AddMap("git ls-files", shexec.Ok(""))
	runner.AddMap("git diff-index", shexec.Ok("foo\n"))
	repo = NewRepo("")
	repo.Dir = DIR_VALID
	NewScript(NewAppContext(WithCommandRunner(runner)), repo).GetIndexStatus()
	if repo.Index != INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", INDEX_MODIFIED, repo.Index)
	}
}
