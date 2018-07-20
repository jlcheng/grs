package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/core"
	"os"
)

type Args struct {
	Clear bool
	Show  bool
}

func main() {
	args := Args{}
	flag.BoolVar(&args.Clear, "clear", false, "clears all notifications for all repos")
	flag.BoolVar(&args.Show, "show", false, "prints the ack time for each repo")
	flag.Parse()

	// TODO: Sets the 'ack_time' of each known repo to the current time
	ctx := grs.NewAppContext()
	sctx, err := grs.InitScriptCtx(config.NewConfigParams(), ctx)
	if err != nil {
		grs.Info("%v", err)
		os.Exit(1)
	}

	repos := sctx.Repos
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		fmt.Printf("create %v if it doesn't exist\n", config.UserConf)
		os.Exit(1)
	}

	if args.Clear {

	}

	if args.Show {
		db, err := ctx.DBService().LoadDB(config.UserDBName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if bytes, err := json.MarshalIndent(db, "", "  "); err == nil {
			fmt.Println(string(bytes))
		}
	}
}
