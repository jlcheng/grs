package main

import (
	"fmt"
	"jcheng/grs/gitscripts"
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
		gitscripts.SetLogLevel(gitscripts.DEBUG)
	}

	repo := defaultRepo(args)
	c, err := gitscripts.Status(repo)
	if err != nil {
		fmt.Println("err ", err)
		return
	}

	fmt.Println(c.String())
}

// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepo(args Args) gitscripts.Repo {
	if len(args.repo) != 0 {
		gitscripts.Debug("Using repo from CLI: %v", args.repo)
		return gitscripts.Repo{args.repo}
	}

	val, ok := os.LookupEnv("GRS_DEFAULT")
	if ok {
		gitscripts.Debug("Using repo from $GRS_DEFAULT: %v", val)
		return gitscripts.Repo{val}
	}

	gitscripts.Debug("Using default repo %v", os.ExpandEnv("$HOME/grstest"))
	return gitscripts.Repo{os.ExpandEnv("$HOME/grstest")}
}