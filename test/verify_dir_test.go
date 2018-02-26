
package test

import (
	"testing"
	"jcheng/grs/grs"
	"os"
	"jcheng/grs/script"
	"jcheng/grs/status"
)

func TestBeforeScript_Fail(t *testing.T) {
	runner := NewMockRunner()
	var repo grs.Repo

	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:cwd}
	}
	rstat := status.NewRStat()
	runner.Add(Error("failed"))
	script.BeforeScript(grs.NewAppContext(), repo, runner, rstat)
	if rstat.Dir == status.DIR_VALID {
		t.Errorf("expected %s, got: %v\n" +
			"", status.DIR_INVALID, rstat.Dir)
	}
}

func TestBeforeScript_OK(t *testing.T) {
	runner := NewMockRunner()
	var repo grs.Repo

	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = grs.Repo{Path:cwd}
	}
	rstat := status.NewRStat()
	runner.Add(Ok(""))
	script.BeforeScript(grs.NewAppContext(), repo, runner, rstat)
	if rstat.Dir != status.DIR_VALID {
		t.Errorf("expected %s, got: %v\n", status.DIR_VALID, rstat.Dir)
	}
}

