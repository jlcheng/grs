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

	// with for the gui goroutine to be ready
	started := make(chan int)
	go func() {
		close(started)
		gui.Start()
	}()
	<-started

	syncController := script.NewSyncController(repos, ctx, displayCh)

	// always run at least once
	syncController.Run()
	if args.daemon {
		ticker := time.NewTicker(time.Duration(args.refresh) * time.Second)
		defer ticker.Stop() // remove? not strictly necessary as we don't offer a way to gracefully shutdown

		// use Ctrl-C to stop this program
		for {
			select {
			case <-ticker.C:
				syncController.Run()
			}
		}
	}

	close(displayCh)
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