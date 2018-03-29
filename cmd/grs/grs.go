package main

import (
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"jcheng/grs/grsdb"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"time"
	"jcheng/grs/display"
)

type Args struct {
	repos   string
	verbose bool
	command string
}

func main() {

	args := Args{}
	flag.StringVar(&args.repos, "repos", "", "target repos")
	flag.StringVar(&args.command, "command", "", "command to run")
	flag.BoolVar(&args.verbose, "verbose", false, "verbosity")
	flag.Parse()

	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	if err := config.SetupUserPrefDir(config.UserPrefDir); err != nil {
		grs.Info("Cannot create user preference directory [%v]:%v", err)
		return
	}

	runner := grs.ExecRunner{}

	ctx := grs.NewAppContext()

	cp := config.NewConfigParams()
	conf, err := config.ReadConfig(cp)
	if conf != nil {
		if conf.Git != "" {
			ctx.SetGitExec(conf.Git)
		}
	} else {
		grs.Debug("configuration error: %v", err)
	}

	repos := grs.ReposFromConf(conf.Repos)
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		fmt.Printf("create %v if it doesn't exist\n", config.UserConf)
		os.Exit(1)
	}

	if db, err := grsdb.LoadFile(ctx.DbPath); err == nil {
		ctx.SetDB(db)
	}


	var screen display.Display = display.NewAnsiDisplay(os.Stdout)

	var repoStatusList = make([]display.RepoStatus, len(repos))
	for idx, repo := range repos {
		rstat := status.NewRStat()
		script.BeforeScript(ctx, repo, runner, rstat)
		if rstat.Dir == status.DIR_VALID {
			script.Fetch(ctx, runner, rstat, repo)
		}
		if rstat.Dir == status.DIR_VALID {
			script.GetRepoStatus(ctx, runner, rstat)

		}
		if rstat.Dir == status.DIR_VALID {
			script.GetIndexStatus(ctx, runner, rstat)
		}

		merged := false
		if atime, err := script.GetActivityTime(repo.Path); err == nil && time.Now().After(atime.Add(ctx.ActivityTimeout)) {
			merged = script.AutoFFMerge(ctx, runner, rstat)
		}

		if repoPtr := ctx.DB().FindRepo(repo.Path); repoPtr != nil {
			repoPtr.RStat.Update(*rstat)
		}

		repoStatusList[idx] = display.RepoStatus{
			Path: repo.Path,
			Rstat: *rstat,
			Merged: merged,
		}
	}
	grsdb.SaveFile(ctx.DBWriter(), ctx.DbPath, ctx.DB())
	screen.SummarizeRepos(repoStatusList)
	screen.Update()
}
