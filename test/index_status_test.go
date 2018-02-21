package test

import (
	"testing"
	"jcheng/grs/status"
	"jcheng/grs/script"
	"errors"
)

func TestGetIndexStatus_Ls_Files_Failed(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", NewCommandHelper([]byte(""), errors.New("failed")))
	runner.AddMap("git diff-index", NewCommandHelper([]byte(""),nil))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, rstat.Index)
	}
}

func TestGetIndexStatus_Diff_Index_Failed(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", NewCommandHelper([]byte(""),nil))
	runner.AddMap("git diff-index", NewCommandHelper([]byte(""), errors.New("failed")))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.INDEX_UNKNOWN, rstat.Index)
	}
}

func TestGetIndexStatus_Unmodified_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git ls-files", NewCommandHelper([]byte(""),nil))
	runner.AddMap("git diff-index", NewCommandHelper([]byte(""),nil))
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
	runner.AddMap("git ls-files", NewCommandHelper([]byte("foo\n"),nil))
	runner.AddMap("git diff-index", NewCommandHelper([]byte(""),nil))
	rstat = status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, rstat.Index)
	}

	runner = NewMockRunner()
	runner.AddMap("git ls-files", NewCommandHelper([]byte(""),nil))
	runner.AddMap("git diff-index", NewCommandHelper([]byte("foo\n"),nil))
	rstat = status.NewRStat()
	rstat.Dir = status.DIR_VALID
	script.GetIndexStatus(runner, rstat)
	if rstat.Index != status.INDEX_MODIFIED {
		t.Errorf("expected %s, got: %v\n", status.INDEX_MODIFIED, rstat.Index)
	}
}