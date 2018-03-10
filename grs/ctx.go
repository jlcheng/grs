package grs

import (
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
	"path/filepath"
	"os"
	"time"
)

type AppContext struct {
	confParams *config.ConfigParams
	cliRepos []string
	defaultGitExec string
	dbWriter grsdb.DBWriter
	db *grsdb.DB
	MinFetchSec int
	ActivityTimeout time.Duration
	DbPath string
}

func NewAppContext() *AppContext {
	return &AppContext{
		confParams: config.NewConfigParams(),
		defaultGitExec: "git",
		cliRepos: []string{},
		dbWriter: grsdb.FileDBWriter,
		db: &grsdb.DB{Repos:make([]grsdb.Repo,0)},
		MinFetchSec: 60 * 60,
		ActivityTimeout: 2 * time.Hour,
		DbPath: filepath.Join(os.ExpandEnv("${HOME}"),".grsdb.json"),
	}
}

func (ctx *AppContext) CliRepos(cliRepos []string) {
	ctx.cliRepos = cliRepos
}

func (ctx *AppContext) ConfParams(confParams *config.ConfigParams) {
	ctx.confParams = confParams
}

func (ctx *AppContext) GetRepos() []string {
	if len(ctx.cliRepos) != 0 {
		return ctx.cliRepos
	}

	if c, err := config.GetCurrConfig(ctx.confParams); err == nil {
		r := make([]string,len(c.Repos))
		for idx,elem := range c.Repos {
			r[idx] = elem.Path
		}
		return r
	}
	return make([]string,0)
}

func (ctx *AppContext) GetGitExec() string {
	if c, err := config.GetCurrConfig(ctx.confParams); err == nil && len(c.Git) != 0  {
		return c.Git
	}

	return ctx.defaultGitExec
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