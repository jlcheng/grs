package main

import (
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"os"
	"jcheng/grs/script"
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

	_ = sctx
	runner := grs.ExecRunner{}
	for idx, elem := range repos {
		_ = idx
		rstat := &status.RStat{}
		script.BeforeScript(ctx, elem, runner, rstat)
		script.AutoRebase(ctx, elem, runner, rstat, false)
	}
}
