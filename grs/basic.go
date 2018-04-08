package grs

import (
	"os/exec"
	"bytes"
	"os"
	"time"
	"errors"
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
)

type Repo struct {
	Path string
}

type Result struct {
	delegate *exec.Cmd
	Stdout string
}

func Status(repo Repo) (*Result, error) {
	cmd := new(Result)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}

	cmd.delegate = exec.Command("git","status")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func Pwd(repo Repo) (*Result, error) {
	cmd := new(Result)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}

	cmd.delegate = exec.Command("pwd")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func Rebase(repo Repo) (*Result, error) {
	cmd := new(Result)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}

	if repo.modifiedRecently() {
		Debug("%v was modified recently. Will not rebase.", repo.Path)
		return cmd, errors.New("repo recently modified")
	}


	buf := new(bytes.Buffer)
	cmd.delegate = exec.Command("git","fetch")
	cmd.delegate.Stdout = buf
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}

	cmd.delegate = exec.Command("git", "rebase", "origin/master")
	cmd.delegate.Stdout = buf
	err = cmd.delegate.Run()
	if err != nil {
		cmd.delegate = exec.Command("git", "rebase", "--abort")
		cmd.delegate.Stdout = buf
		err = cmd.delegate.Run()
		if err != nil {
			return cmd, err
		}
		Debug("Unable to rebase. Original HEAD restored.")
		return cmd, nil
	}
	return cmd, nil
}

func RepoPath(repo Repo) (*Result, error) {
	cmd := exec.Command("")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	buf.WriteString(repo.Path)
	result := Result{delegate:cmd}
	return &result, nil
}

func (repo *Repo) modifiedRecently() bool {
	info, err := os.Stat(repo.Path)
	if err != nil {
		Debug("%v stat failed, modification time unknown", repo.Path)
		return true
	}
	tdiff := time.Now().Sub(info.ModTime())
	Debug("%v last accessed %v ago", repo.Path, tdiff)
	if tdiff > time.Hour {
		return false
	} else {
		return true
	}
}

func (cmd *Result) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}

func ReposFromConf(rc []config.RepoConf) []Repo {
	var r = make([]Repo,len(rc))
	for idx, elem := range rc {
		r[idx] = Repo{Path:elem.Path}
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


	// initialize ~/.grs.d/ kvstore
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