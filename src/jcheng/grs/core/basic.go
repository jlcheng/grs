package grs

import (
	"bytes"
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
	"jcheng/grs/status"
	"os"
	"os/exec"
	"strings"
)

type Result struct {
	delegate *exec.Cmd
	Stdout   string
}

func (cmd *Result) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}

func ReposFromConf(rc []config.RepoConf) []status.Repo {
	var r = make([]status.Repo, len(rc))
	for idx, elem := range rc {
		r[idx] = status.Repo{Path: elem.Path}
	}
	return r
}

func ReposFromString(input string) []status.Repo {
	tokens := strings.Split(input, string(os.PathListSeparator))
	r := make([]status.Repo, len(tokens))
	for idx, elem := range tokens {
		r[idx] = status.Repo{Path: elem}
	}
	return r
}

func InitScriptCtx(cp *config.ConfigParams, ctx *AppContext) (*ScriptContext, error) {
	if err := config.SetupUserPrefDir(config.UserPrefDir); err != nil {
		return nil, err
	}

	// read ~/.grs.d/config.json
	sctx := NewScriptContext(ctx)
	conf, err := config.ReadConfig(cp)
	if conf != nil {
		if conf.Git != "" {
			ctx.SetGitExec(conf.Git)
		}
	} else {
		return nil, err
	}
	sctx.Repos = ReposFromConf(conf.Repos)

	// initialize ~/.grs.d kvstore
	if kvstore, err := grsdb.InitDiskKVStore(config.UserPrefDir); err == nil {
		ctx.dbService = grsdb.NewDBService(kvstore)
	} else {
		return nil, err
	}

	// read ~/.grs.d/grs.db
	if db, err := ctx.DBService().LoadDB(config.UserDBName); err == nil {
		ctx.SetDB(db)
	} else if os.IsNotExist(err) {
		ctx.SetDB(&grsdb.DB{})
	}

	return sctx, nil
}