package test

import (
	"jcheng/grs/config"
	"jcheng/grs/core"
	"reflect"
	"testing"
)

// TestGetRepos_ConfFile verifies resolving repos from ConfigParams
func TestGetRepos_ConfFile(t *testing.T) {
	var conf *config.Config
	cp := &config.ConfigParams{Env: "data/config.json", User: "data/empty_config.json"}
	conf, _ = config.ReadConfig(cp)

	if r := RepoIds(conf.Repos); !reflect.DeepEqual([]string{"rel/repo1", "/abs/repo2"}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}

	cp = &config.ConfigParams{User: "data/config.json"}
	conf, _ = config.ReadConfig(cp)
	if r := RepoIds(conf.Repos); !reflect.DeepEqual([]string{"rel/repo1", "/abs/repo2"}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}

	cp = &config.ConfigParams{}
	conf, _ = config.ReadConfig(cp)
	if r := RepoIds(conf.Repos); !reflect.DeepEqual([]string{}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}
}

// TestGetGitExec_ConfFile verifies that GetGitExec() is controlled by ConfigParams
func TestGetGitExec_ConfFile(t *testing.T) {
	ctx := grs.NewAppContext()
	cp := &config.ConfigParams{User: "data/config.json"}
	if conf, _ := config.ReadConfig(cp); conf != nil {
		ctx.SetGitExec(conf.Git)
	}

	if r := ctx.GetGitExec(); r != "/path/to/git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}

// Verifies that the default GetGitExec() is `git`
func TestGetGitExecDefault(t *testing.T) {
	ctx := grs.NewAppContext()

	if r := ctx.GetGitExec(); r != "git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}

func RepoIds(repos []config.RepoConf) []string {
	retval := make([]string, len(repos))
	for i, repo := range repos {
		retval[i] = repo.Path
	}
	return retval
}
