package script

import (
	"errors"
	"jcheng/grs/config"
	"jcheng/grs/core"
	"jcheng/grs/grsdb"
	"os"
	"path/filepath"
	"time"
)

var lastActivityFiles = []string{"HEAD", "COMMIT_EDITMSG", "ORIG_HEAD", "index", "config"}

// GetActivityTime gets the estimated "last modified time" of a repo
func GetActivityTime(repo string) (time.Time, error) {
	var atime time.Time
	if f, err := os.Stat(repo); err != nil || !f.IsDir() {
		return atime, errors.New("%v is not a directory")
	}
	for _, f := range lastActivityFiles {
		fn := filepath.Join(repo, ".git", f)
		if finfo, err := os.Stat(fn); err == nil {
			if finfo.ModTime().After(atime) {
				atime = finfo.ModTime()
			}
		}
	}
	return atime, nil
}

func InitScriptCtx(cp *config.ConfigParams, ctx *grs.AppContext) (*grs.ScriptContext, error) {
	if err := config.SetupUserPrefDir(config.UserPrefDir); err != nil {
		return nil, err
	}

	// read ~/.grs.d/config.json
	sctx := grs.NewScriptContext(ctx)
	conf, err := config.ReadConfig(cp)
	if conf != nil {
		if conf.Git != "" {
			ctx.SetGitExec(conf.Git)
		}
	} else {
		return nil, err
	}
	sctx.Repos = grs.ReposFromConf(conf.Repos)

	// initialize ~/.grs.d kvstore
	if kvstore, err := grsdb.InitDiskKVStore(config.UserPrefDir); err == nil {
		ctx.DbService = grsdb.NewDBService(kvstore)
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
