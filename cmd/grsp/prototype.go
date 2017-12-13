package main

import (
	"fmt"
	"jcheng/grs/grs"
	"os"
	"flag"
)

type Args struct {
	repo string
	verbose bool
}

func main() {
	flag.Parse()
	args := Args{repo: flag.Arg(0), verbose: true}
	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	repo := defaultRepo(args)

	var cmd grs.Cmd = grs.Status

	c, err := cmd(repo)
	if err != nil {
		fmt.Println("err ", err)
		return
	}

	fmt.Println(c.String())
}

// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepo(args Args) grs.Repo {
	if len(args.repo) != 0 {
		grs.Debug("Using repo from CLI: %v", args.repo)
		return grs.Repo{args.repo}
	}

	val, ok := os.LookupEnv("GRS_DEFAULT")
	if ok {
		grs.Debug("Using repo from $GRS_DEFAULT: %v", val)
		return grs.Repo{val}
	}

	grs.Debug("Using default repo %v", os.ExpandEnv("$HOME/grstest"))
	return grs.Repo{os.ExpandEnv("$HOME/grstest")}
}