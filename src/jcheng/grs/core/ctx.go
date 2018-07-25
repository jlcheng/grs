package grs

import (
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
	"jcheng/grs/status"
	"time"
)

type AppContext struct {
	CommandRunner
	confParams      *config.ConfigParams
	defaultGitExec  string
	db              *grsdb.DB
	MinFetchSec     int
	ActivityTimeout time.Duration
	DbPath          string
	DbService       grsdb.DBService
}

func NewAppContext() *AppContext {
	return &AppContext{
		confParams:      config.NewConfigParams(),
		defaultGitExec:  "git",
		db:              &grsdb.DB{Repos: make([]grsdb.RepoDTO, 0)},
		MinFetchSec:     60 * 60,
		ActivityTimeout: 2 * time.Hour,
		DbPath:          config.UserDB,
	}
}
func NewAppContextWithRunner(runner CommandRunner) *AppContext {
	return &AppContext{
		CommandRunner:   runner,
		confParams:      config.NewConfigParams(),
		defaultGitExec:  "git",
		db:              &grsdb.DB{Repos: make([]grsdb.RepoDTO, 0)},
		MinFetchSec:     60 * 60,
		ActivityTimeout: 2 * time.Hour,
		DbPath:          config.UserDB,
	}
}

func (ctx *AppContext) ConfParams(confParams *config.ConfigParams) {
	ctx.confParams = confParams
}

func (ctx *AppContext) GetGitExec() string {
	return ctx.defaultGitExec
}

func (ctx *AppContext) SetGitExec(defaultGitExec string) {
	ctx.defaultGitExec = defaultGitExec
}

func (ctx *AppContext) DB() *grsdb.DB {
	return ctx.db
}

func (ctx *AppContext) SetDB(db *grsdb.DB) {
	ctx.db = db
}

func (ctx *AppContext) DBService() grsdb.DBService {
	return ctx.DbService
}

type ScriptContext struct {
	Ctx   *AppContext
	Repos []status.Repo
}

func NewScriptContext(ctx *AppContext) *ScriptContext {
	return &ScriptContext{
		Ctx:   ctx,
		Repos: make([]status.Repo, 0),
	}
}
