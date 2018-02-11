package test

import (
	"testing"
	"jcheng/grs/grs"
	"reflect"
	"jcheng/grs/config"
)

func TestGetReposCli(t *testing.T) {
	ctx := grs.GetContext()
	in := []string{"cli/rel/repo1","/cli/abs/repo2"}
	ctx.CliRepos(in)

	if r := ctx.GetRepos(); !reflect.DeepEqual(in, r) {
		t.Error("Unexpected repos. Got:", r)
	}
}

func TestGetReposConfFile(t *testing.T) {
	ctx := grs.GetContext()
	ctx.CliRepos([]string{})

	ctx.ConfParams(&config.ConfigParams{Env: "data/config.json", User: "data/empty_config.json"})
	if r := ctx.GetRepos(); !reflect.DeepEqual([]string{"rel/repo1","/abs/repo2"}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}

	ctx.ConfParams(&config.ConfigParams{User: "data/config.json"})
	if r := ctx.GetRepos(); !reflect.DeepEqual([]string{"rel/repo1","/abs/repo2"}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}

	ctx.ConfParams(&config.ConfigParams{})
	if r := ctx.GetRepos(); !reflect.DeepEqual([]string{}, r) {
		t.Error("Unexpected repos. Got: ", r)
	}
}

func TestGetReposCliAndConfFile(t *testing.T) {
	ctx := grs.GetContext()
	in := []string{"cli/rel/repo1","/cli/abs/repo2"}
	ctx.CliRepos(in)
	ctx.ConfParams(&config.ConfigParams{User: "data/config.json"})

	if r := ctx.GetRepos(); !reflect.DeepEqual(in, r) {
		t.Error("Unexpected repos. Got:", r)
	}
}

func TestGetGitExecConfFile(t *testing.T) {
	ctx := grs.GetContext()
	ctx.ConfParams(&config.ConfigParams{User: "data/config.json"})

	if r := ctx.GetGitExec(); r != "/path/to/git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}

// Verifies that the default git exec is `git`
func TestGetGitExecDefault(t *testing.T) {
	ctx := grs.GetContext().ResetInternal()

	if r := ctx.GetGitExec(); r != "git" {
		t.Error("Unexpected git executable. Got:", r)
	}
}
