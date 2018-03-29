package grs

import (
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
	"path/filepath"
	"time"
)

type AppContext struct {
	confParams      *config.ConfigParams
	cliRepos        []string
	defaultGitExec  string
	dbWriter        grsdb.DBWriter
	db              *grsdb.DB
	MinFetchSec     int
	ActivityTimeout time.Duration
	DbPath          string
}

func NewAppContext() *AppContext {
	return &AppContext{
		confParams:      config.NewConfigParams(),
		defaultGitExec:  "git",
		cliRepos:        []string{},
		dbWriter:        grsdb.FileDBWriter,
		db:              &grsdb.DB{Repos: make([]grsdb.Repo, 0)},
		MinFetchSec:     60 * 60,
		ActivityTimeout: 2 * time.Hour,
		DbPath:          filepath.Join(config.UserDB),
	}
}

func (ctx *AppContext) CliRepos(cliRepos []string) {
	ctx.cliRepos = cliRepos
}

func (ctx *AppContext) ConfParams(confParams *config.ConfigParams) {
	ctx.confParams = confParams
}

func (ctx *AppContext) GetRepos() []string {
	return ctx.cliRepos
}

func (ctx *AppContext) SetRepos(repos []string) {
	ctx.cliRepos = repos
}

func (ctx *AppContext) GetGitExec() string {
	return ctx.defaultGitExec
}

func (ctx *AppContext) SetGitExec(defaultGitExec string) {
	ctx.defaultGitExec = defaultGitExec
}

func (ctx *AppContext) DBWriter() grsdb.DBWriter {
	return ctx.dbWriter
}

func (ctx *AppContext) DB() *grsdb.DB {
	return ctx.db
}

func (ctx *AppContext) SetDB(db *grsdb.DB) {
	ctx.db = db
}