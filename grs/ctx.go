package grs

import (
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
	"time"
)

type AppContext struct {
	confParams      *config.ConfigParams
	defaultGitExec  string
	db              *grsdb.DB
	MinFetchSec     int
	ActivityTimeout time.Duration
	DbPath          string
	dbService       grsdb.DBService
}

func NewAppContext() *AppContext {
	return &AppContext{
		confParams:      config.NewConfigParams(),
		defaultGitExec:  "git",
		db:              &grsdb.DB{Repos: make([]grsdb.Repo, 0)},
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
	return ctx.dbService
}

type ScriptContext struct {
	Ctx *AppContext
	Repos []Repo
}

func NewScriptContext(ctx *AppContext) *ScriptContext {
	return &ScriptContext{
		Ctx: ctx,
		Repos: make([]Repo, 0),
	}
}