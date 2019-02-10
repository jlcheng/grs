package script

import (
	"jcheng/grs/shexec"
	"os"
	"testing"
)

func TestBeforeScript_Fail(t *testing.T) {
	runner := shexec.NewMockRunner()
	var repo *Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = NewRepo(cwd)
	}
	runner.Add(shexec.Error("failed"))
	s := NewScript(shexec.NewAppContext(shexec.WithCommandRunner(runner)), repo)
	s.BeforeScript()
	if repo.Dir == DIR_VALID {
		t.Errorf("expected %s, got: %v\n"+
			"", DIR_INVALID, repo.Dir)
	}
}

func TestBeforeScript_OK(t *testing.T) {
	runner := shexec.NewMockRunner()
	var repo *Repo
	if cwd, e := os.Getwd(); e != nil {
		t.Error(e)
	} else {
		repo = NewRepo(cwd)
	}

	runner.Add(shexec.Ok(""))
	s := NewScript(shexec.NewAppContext(shexec.WithCommandRunner(runner)), repo)
	s.BeforeScript()
	if repo.Dir != DIR_VALID {
		t.Errorf("expected %s, got: %v\n", DIR_VALID, repo.Dir)
	}
}
