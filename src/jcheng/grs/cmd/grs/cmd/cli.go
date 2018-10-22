package cmd

import (
	"fmt"
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"time"
)

type Args struct {
	repos      []string
	verbose    bool
	command    string
	daemon     bool
	refresh    int
	forceMerge bool
}

func RunCli(args Args) {
	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	ctx := grs.NewAppContextWithRunner(&grs.ExecRunner{})
	repos := grs.ReposFromStringSlice(args.repos)
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	displayCh := make(chan bool)
	var reporter = ReporterStruct{
		ctx: ctx,
		repos: repos,
	}
	gui := script.NewGUI(ctx, displayCh, reporter.Report, args.daemon)
	gui.Start()

	syncDaemon := script.NewSyncDaemon(repos, ctx, displayCh)
	syncDaemon.StartDaemon()

	// always run at least once
	syncDaemon.Run()
	if args.daemon {
		ticker := time.NewTicker(time.Duration(args.refresh) * time.Second)
		defer ticker.Stop() // remove? not strictly necessary as we don't offer a way to gracefully shutdown

		// use CTRL-C to stop this loop
		for true {
			select {
			case <-ticker.C:
				syncDaemon.Run()
			}
		}
	}

	syncDaemon.Shutdown()
	syncDaemon.WaitForShutdown()

	gui.Shutdown()
	gui.WaitShutdown()
}

// TODO: Move to script directory
type ReporterStruct struct {
	ctx *grs.AppContext
	repos []status.Repo
}

func (rs *ReporterStruct) Report() []status.Repo {
	for idx, _ := range rs.repos {
		s := script.NewScript(rs.ctx, &rs.repos[idx])
		s.BeforeScript()
		s.GetRepoStatus()
		s.GetIndexStatus()
	}
	return rs.repos
}