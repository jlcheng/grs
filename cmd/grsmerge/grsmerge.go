package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"os"
	"path/filepath"
	"strings"
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
