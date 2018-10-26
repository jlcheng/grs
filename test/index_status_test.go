package test

import (
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"testing"
)

func TestGetIndexStatus_Ls_Files_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Error("failed"))
	runner.AddMap("git diff-index", Ok(""))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != script.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", script.INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Diff_Index_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Error("failed"))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != script.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", script.INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Unmodified_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok(""))
	repo := script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != script.INDEX_UNMODIFIED {
		t.Errorf("expected %s, got: %v\n", script.INDEX_UNMODIFIED, repo.Index)
	}
}

func TestGetIndexStatus_Modified_Ok(t *testing.T) {
	var runner *MockRunner
	var repo *script.Repo
	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok("foo\n"))
	runner.AddMap("git diff-index", Ok(""))
	repo = script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != script.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", script.INDEX_MODIFIED, repo.Index)
	}

	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok("foo\n"))
	repo = script.NewRepo("")
	repo.Dir = script.DIR_VALID
	script.NewScript(shexec.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != script.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", script.INDEX_MODIFIED, repo.Index)
	}
}
