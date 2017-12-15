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

	repo := defaultRepo(args)

	cmd := defaultCommand(args)

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

	val, ok := os.LookupEnv("GRS_REPO")
	if ok {
		grs.Debug("Using repo from $GRS_REPO: %v", val)
		return grs.Repo{val}
	}

	grs.Debug("Using default repo %v", os.ExpandEnv("$HOME" + string(os.PathSeparator) + "grstest"))
	return grs.Repo{os.ExpandEnv("$HOME/grstest")}
}

func defaultCommand(args Args) grs.Cmd {
	var val string
	if len(args.command) != 0 {
		val = args.command
		grs.Debug("Using action from CLI: %v", args.command)
	}

	if len(val) == 0 {
		var ok bool
		val, ok = os.LookupEnv("GRS_ACTION")
		if ok {
			grs.Debug("Using repo from $GRS_ACTION: %v", val)
		}
	}

	if len(val) == 0 {
		grs.Debug("Using default action: status")
		val = "status"
	}

	switch val {
	case "pwd":
		return grs.Pwd
	case "rebase":
		return grs.Rebase
	default:
		return grs.Status
	}
}