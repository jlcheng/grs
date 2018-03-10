package main

import (
	"jcheng/grs/grs"
	"os"
	"flag"
	"fmt"
	"jcheng/grs/script"
	"strings"
	"jcheng/grs/status"
	"jcheng/grs/config"
	"jcheng/grs/grsdb"
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

	runner := grs.ExecRunner{}
	repos := defaultRepos(args)
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	ctx := grs.NewAppContext()
	if db, err := grsdb.LoadFile(ctx.DbPath); err == nil {
		ctx.SetDB(db)
	}
	for _, repo := range repos {
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
		if atime, err := script.GetActivityTime(repo.Path);
			err == nil && time.Now().After(atime.Add(ctx.ActivityTimeout)) {
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



// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepos(args Args) []grs.Repo {
	if len(args.repos) != 0 {
		grs.Debug("Using repos from CLI: %v", args.repos)
		return mkrepos(args.repos)
	}

	p := config.NewConfigParams()
	if c, _ := config.GetCurrConfig(p); c != nil {
		return grs.ReposFromConf(c.Repos)
	}
	return []grs.Repo{}
}

func mkrepos(s string) []grs.Repo {
	res := make([]grs.Repo, 0, 1)
	for _, elem := range strings.Split(s, string(os.PathListSeparator)) {
		res = append(res, grs.Repo{Path:elem})
	}
	return res
}
