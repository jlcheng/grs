package main

import (
	"jcheng/grs/grs"
	"os"
	"flag"
	"fmt"
	"jcheng/grs/script"
)

type Args struct {
	repo string
	verbose bool
	command string
}

func main() {

	args := Args{}
	flag.StringVar(&args.repo, "repo", "", "target repo")
	flag.StringVar(&args.command, "command", "", "command to run")
	flag.BoolVar(&args.verbose, "verbose", true, "verbosity")
	flag.Parse()

	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	runner := grs.ExecRunner{}
	repo := defaultRepo(args)
	script := defaultScript(args)
	fmt.Printf("repo [%v] status is %v\n", repo.Path, script(repo, runner))
}

// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepo(args Args) grs.Repo {
	if len(args.repo) != 0 {
		grs.Debug("Using repo from CLI: %v", args.repo)
		return grs.Repo{args.repo}
	}

	val, ok := os.LookupEnv("GRS_REPO")
	if ok {
		grs.Debug("Using repo from $GRS_REPO: %v", val)
		return grs.Repo{val}
	}

	grs.Debug("Using default repo %v", os.ExpandEnv("$HOME" + string(os.PathSeparator) + "grstest"))
	return grs.Repo{os.ExpandEnv("$HOME/grstest")}
}

func defaultScript(args Args) script.Script {
	grs.Debug("Using hard-coded script `GetRepoStatus`")
	return script.GetRepoStatus
}