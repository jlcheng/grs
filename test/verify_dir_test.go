package test

import (
	"jcheng/grs/grs"
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
	script.BeforeScript(grs.NewAppContextWithRunner(runner), repo)
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
	script.BeforeScript(grs.NewAppContextWithRunner(runner), repo)
	if repo.Dir != status.DIR_VALID {
		t.Errorf("expected %s, got: %v\n", status.DIR_VALID, repo.Dir)
	}
}
