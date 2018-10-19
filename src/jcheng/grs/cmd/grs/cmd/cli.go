package cmd

import (
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/core"
	"os"
)

type Args struct {
	repos      string
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

	display := make(chan bool)
	syncDaemon := grs.NewSyncDaemon(repos, ctx, display)
	syncDaemon.StartDaemon()
	syncDaemon.Shutdown()
	syncDaemon.WaitForShutdown()

}
