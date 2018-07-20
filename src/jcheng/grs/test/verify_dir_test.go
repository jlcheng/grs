package test

import (
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"testing"
)

func TestBeforeScript_Fail(t *testing.T) {
	runner := NewMockRunner()
	var repo *status.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = status.NewRepo(cwd)
	}
	runner.Add(Error("failed"))
	s := script.NewScript(grs.NewAppContextWithRunner(runner), repo)
	s.BeforeScript()
	if repo.Dir == status.DIR_VALID {
		t.Errorf("expected %s, got: %v\n"+
			"", status.DIR_INVALID, repo.Dir)
	}
}

func TestBeforeScript_OK(t *testing.T) {
	runner := NewMockRunner()
	var repo *status.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = status.NewRepo(cwd)
	}

	runner.Add(Ok(""))
	s := script.NewScript(grs.NewAppContextWithRunner(runner), repo)
	s.BeforeScript()
	if repo.Dir != status.DIR_VALID {
		t.Errorf("expected %s, got: %v\n", status.DIR_VALID, repo.Dir)
	}
}
