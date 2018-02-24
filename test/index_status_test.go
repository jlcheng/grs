package test

import (
	"testing"
	"jcheng/grs/status"
	"jcheng/grs/script"
)

func TestGetIndexStatus_Ls_Files_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Error("failed"))
	runner.AddMap("git diff-index", Ok(""))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, rstat.Index)
	}
}

func TestGetIndexStatus_Diff_Index_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Error("failed"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, rstat.Index)
	}
}

func TestGetIndexStatus_Unmodified_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok(""))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_UNMODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNMODIFIED, rstat.Index)
	}
}

func TestGetIndexStatus_Modified_Ok(t *testing.T) {
	var runner *MockRunner
	var rstat *status.RStat
	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok("foo\n"))
	runner.AddMap("git diff-index", Ok(""))
	rstat = status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, rstat.Index)
	}

	runner = NewMockRunner()
	runner.AddMap("git ls-files", Ok(""))
	runner.AddMap("git diff-index", Ok("foo\n"))
	rstat = status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, rstat.Index)
	}
}