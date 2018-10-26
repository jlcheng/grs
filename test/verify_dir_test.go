package test

import (
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"os"
	"testing"
)

func TestBeforeScript_Fail(t *testing.T) {
	runner := NewMockRunner()
	var repo *script.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = script.NewRepo(cwd)
	}
	runner.Add(Error("failed"))
	s := script.NewScript(shexec.NewAppContextWithRunner(runner), repo)
	s.BeforeScript()
	if repo.Dir == script.DIR_VALID {
		t.Errorf("expected %s, got: %v\n"+
			"", script.DIR_INVALID, repo.Dir)
	}
}

func TestBeforeScript_OK(t *testing.T) {
	runner := NewMockRunner()
	var repo *script.Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = script.NewRepo(cwd)
	}

	runner.Add(Ok(""))
	s := script.NewScript(shexec.NewAppContextWithRunner(runner), repo)
	s.BeforeScript()
	if repo.Dir != script.DIR_VALID {
		t.Errorf("expected %s, got: %v\n", script.DIR_VALID, repo.Dir)
	}
}
