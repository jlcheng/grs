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
	// Set up a working directory and update some sort of metadata object (grsdb)
	// for any repo that requires rebasing (branch == diverged):
	//  1. Set up a clone directory
	//  2. Identify merge-base
	//  3. Identify the graph of child commits from merge-base to HEAD
	//  4. Abort if a node with multiple parents is identified; otherwise call it lineage
	//  5. Rebase current branch against each child in lineage
	//  6. Stop when conflict is detected
	//  7. Now the cloned branch is "as caught up as possible" against ${upstream}

}
