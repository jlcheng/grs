package main

import (
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/display"
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"io"
	"net/http"
	"os"
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


	ctx := grs.NewAppContextWithRunner(&grs.ExecRunner{})
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

	go UpdateRepos(args, ctx, repos)
	Tmp()
}

func UpdateReposRepeat(args Args, ctx *grs.AppContext, repos []status.Repo) {
	for true {
		UpdateRepos(args, ctx, repos)
		if !args.daemon {
			break
		}
		time.Sleep(time.Second * time.Duration(args.refresh))
	}
}

func UpdateRepos(args Args, ctx *grs.AppContext, repos []status.Repo) {
	screen := display.NewAnsiDisplay(args.daemon, os.Stdout)
	repoStatusList := make([]display.RepoVO, len(repos))
	for idx, repo := range repos {
		s := script.NewScript(ctx, &repo)
		s.BeforeScript()
		s.Fetch()
		s.GetRepoStatus()
		s.GetIndexStatus()

		// allow forced-merge in non-daemon mode. otherwise, use last modified time to decide mergeness
		merged := false
		doMerge := args.force_merge && !args.daemon
		if !doMerge {
			atime, err := script.GetActivityTime(repo.Path)
			doMerge = (err == nil) && time.Now().After(atime.Add(ctx.ActivityTimeout))
		}
		if doMerge {
			switch repo.Branch {
			case status.BRANCH_BEHIND:
				s.AutoFFMerge()
			case status.BRANCH_DIVERGED:
				s.AutoRebase()
			}
			s.GetRepoStatus()
			s.GetIndexStatus()
		}

		repoPtr := ctx.DB().FindOrCreateRepo(repo.Path)
		if repoPtr != nil {
			repoPtr.RStat.Update(&repo)
			if merged {
				repoPtr.MergedCnt = repoPtr.MergedCnt + 1
				repoPtr.MergedSec = time.Now().Unix()
			}
			repoStatusList[idx] = display.RepoVO{
				Repo:      repo,
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
}

func Tmp() {
	http.HandleFunc("/st", func(w http.ResponseWriter, req *http.Request){
		io.WriteString(w, "status\n")
	})
	http.ListenAndServe(":8080", nil)

}
