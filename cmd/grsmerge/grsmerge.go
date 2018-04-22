package main

import (
	"flag"
	"jcheng/grs/grs"
	"jcheng/grs/config"
	"os"
	"fmt"
	"encoding/json"
	"jcheng/grs/status"
)

type Args struct {
	verbose bool
	repos   string
}

func main() {
	args := Args{}
	flag.BoolVar(&args.verbose, "verbose", false, "verbosity")
	flag.StringVar(&args.repos, "repos", "", "target repos")
	flag.Parse()

	ctx := grs.NewAppContext()
	sctx, err := grs.InitScriptCtx(config.NewConfigParams(), ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	repos := grs.ReposFromString(args.repos)
	if repos[0].Path == "" {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	for idx, elem := range repos {
		_ = idx
		_ = elem
	}
}

func AutoRebase(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat, repo grs.Repo) {
	git := ctx.GetGitExec()
	r := *runner.Command(git, "")
	_ = r
	// set up a working directory and update some sort of repo_metadata object
}
