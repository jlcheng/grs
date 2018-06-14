package main

import (
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/display"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"os/signal"
	"time"
)

type Args struct {
	repos       string
	verbose     bool
	command     string
	daemon      bool
	refresh     int
	force_merge bool
}

func main() {

	args := Args{}
	flag.StringVar(&args.repos, "repos", "", "target repos")
	flag.StringVar(&args.command, "command", "", "command to run")
	flag.BoolVar(&args.verbose, "verbose", false, "verbosity")
	flag.BoolVar(&args.daemon, "d", false, "[daemon mode] enable daemon mode")
	flag.IntVar(&args.refresh, "r", 300, "[daemon mode] How often to check for changes, in seconds.")
	flag.BoolVar(&args.force_merge, "merge", false, "ignore access time check when auto-merging")
	flag.Parse()

	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	ctx := grs.NewAppContext()
	sctx, err := grs.InitScriptCtx(config.NewConfigParams(), ctx)
	if err != nil {
		grs.Info("%v", err)
		os.Exit(1)
	}

	repos := grs.ReposFromString(args.repos)
	if repos[0].Path == "" {
		repos = sctx.Repos
	}
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		fmt.Printf("create %v if it doesn't exist\n", config.UserConf)
		os.Exit(1)
	}

	var screen display.Display = display.NewAnsiDisplay(os.Stdout)
	var repoStatusList = make([]display.RepoVO, len(repos))

	ctrl := make(chan os.Signal, 1)
	signal.Notify(ctrl, os.Interrupt)
	go func() {
		for sig := range ctrl {
			grs.Debug("got %v, quitting", sig)
			os.Exit(0)
		}
	}()

	runner := &grs.ExecRunner{}
	var repeat = true
	for repeat {
		for idx, repo := range repos {
			repoTwo := status.NewRepo(repo.Path)
			script.BeforeScript(ctx, runner, repoTwo)
			if repoTwo.Dir != status.DIR_VALID {
				script.Fetch(ctx, runner, repoTwo)
			}
			if repoTwo.Dir == status.DIR_VALID {
				script.GetRepoStatus(ctx, runner, repoTwo)
			}
			if repoTwo.Dir == status.DIR_VALID {
				script.GetIndexStatus(ctx, runner, repoTwo)
			}

			merged := false
			do_merge := false

			// check for recency when in daemon mode, allow forced merge in non-deamon mode
			if !args.daemon && args.force_merge {
				do_merge = true
			}
			if args.daemon || !args.force_merge {
				atime, err := script.GetActivityTime(repoTwo.Path)
				do_merge = (err == nil) && time.Now().After(atime.Add(ctx.ActivityTimeout))
			}
			if repoTwo.Branch != status.BRANCH_UNTRACKED && do_merge {
				merged = script.AutoFFMerge(ctx, runner, repoTwo)
			}

			repoPtr := ctx.DB().FindOrCreateRepo(repo.Path)
			if repoPtr != nil {
				repoPtr.RStat.Update(repoTwo)
				if merged {
					repoPtr.MergedCnt = repoPtr.MergedCnt + 1
					repoPtr.MergedSec = time.Now().Unix()
				}
				repoStatusList[idx] = display.RepoVO{
					Repo:      *repoTwo,
					Merged:    merged,
					MergedSec: repoPtr.MergedSec,
				}
			}
		}
		err := ctx.DBService().SaveDB(config.UserDBName, ctx.DB())
		if err != nil {
			grs.Info("cannot save db %v", err)
		}
		screen.SummarizeRepos(repoStatusList)
		screen.Update()

		if !args.daemon {
			repeat = false
			continue
		}
		time.Sleep(time.Second * time.Duration(args.refresh))
	}
}
