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
	repos := ctx.GetRepos()
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	if db, err := grsdb.LoadFile(ctx.DbPath); err == nil {
		ctx.SetDB(db)
	}

	for _, repoId := range repos {
		repo := grs.Repo{Path: repoId}
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

		if merged {
			grs.Info("repo [%v] auto fast-foward to latest", repo.Path)
		} else {
			grs.Info("repo [%v] status is %v, %v", repo.Path, rstat.Branch, rstat.Index)
		}
	}
	grsdb.SaveFile(ctx.DBWriter(), ctx.DbPath, ctx.DB())
}
