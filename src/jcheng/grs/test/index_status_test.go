package test

import (
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
)

func TestGetIndexStatus_Ls_Files_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Error("failed"))
	runner.AddMap("git diff-index", Ok(""))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.NewScript(grs.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Diff_Index_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Error("failed"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.NewScript(grs.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, repo.Index)
	}
}

func TestGetIndexStatus_Unmodified_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok(""))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.NewScript(grs.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != status.INDEX_UNMODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNMODIFIED, repo.Index)
	}
}

func TestGetIndexStatus_Modified_Ok(t *testing.T) {
	var runner *MockRunner
	var repo *status.Repo
	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok("foo\n"))
	runner.AddMap("git diff-index", Ok(""))
	repo = status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.NewScript(grs.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, repo.Index)
	}

	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok("foo\n"))
	repo = status.NewRepo("")
	repo.Dir = status.DIR_VALID
	script.NewScript(grs.NewAppContextWithRunner(runner), repo).GetIndexStatus()
	if repo.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, repo.Index)
	}
}
