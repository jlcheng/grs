package test

import (
	"testing"
	"jcheng/grs/grs"
	"reflect"
	"jcheng/grs/config"
)

// TestGetReposCli verifies the getter/setter for CliRepos
func TestGetRepos_Cli(t *testing.T) {
	ctx := grs.NewAppContext()
	in := []string{"cli/rel/repo1","/cli/abs/repo2"}
	ctx.CliRepos(in)

	if r := ctx.GetRepos(); !reflect.DeepEqual(in, r) {
		t.Error("Unexpected repos. Got:", r)
	}
}

// TestGetRepos_ConfFile verifies resolving repos from ConfigParams
func TestGetRepos_ConfFile(t *testing.T) {
	ctx := grs.NewAppContext()
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

// TestGetRepos_Cli_And_ConfFile verifies that CLI takes precedence
func TestGetRepos_Cli_And_ConfFile(t *testing.T) {
	ctx := grs.NewAppContext()
	in := []string{"cli/rel/repo1","/cli/abs/repo2"}
	ctx.CliRepos(in)
	ctx.ConfParams(&config.ConfigParams{User: "data/config.json"})

	if r := ctx.GetRepos(); !reflect.DeepEqual(in, r) {
		t.Error("Unexpected repos. Got:", r)
	}
}

// TestGetGitExec_ConfFile verifies that GetGitExec() is controlled by ConfigParams
func TestGetGitExec_ConfFile(t *testing.T) {
	ctx := grs.NewAppContext()
	ctx.ConfParams(&config.ConfigParams{User: "data/config.json"})

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
