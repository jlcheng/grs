package grs

import (
	"jcheng/grs/status"
	"time"
)

type AppContext struct {
	CommandRunner
	defaultGitExec  string
	MinFetchSec     int
	ActivityTimeout time.Duration
	DbPath          string
}

func NewAppContext() *AppContext {
	return &AppContext{
		defaultGitExec:  "git",
		MinFetchSec:     60 * 60,
		ActivityTimeout: 2 * time.Hour,
	}
}
func NewAppContextWithRunner(runner CommandRunner) *AppContext {
	return &AppContext{
		CommandRunner:   runner,
		defaultGitExec:  "git",
		MinFetchSec:     60 * 60,
		ActivityTimeout: 2 * time.Hour,
	}
}

func (ctx *AppContext) GetGitExec() string {
	return ctx.defaultGitExec
}

func (ctx *AppContext) SetGitExec(defaultGitExec string) {
	ctx.defaultGitExec = defaultGitExec
}

type ScriptContext struct {
	Ctx   *AppContext
	Repos []status.Repo
}
