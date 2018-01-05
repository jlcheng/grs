package main

import (
	"jcheng/grs/grs"
	"os"
	"flag"
	"fmt"
	"jcheng/grs/script"
	"strings"
	"jcheng/grs/status"
)

type Args struct {
	repos   string
	verbose bool
	command string
}

func main() {

	args := Args{}
	flag.StringVar(&args.repos, "repos", "", "target repos")
	flag.StringVar(&args.command, "command", "", "command to run")
	flag.BoolVar(&args.verbose, "verbose", true, "verbosity")
	flag.Parse()

	if args.verbose {
		grs.SetLogLevel(grs.DEBUG)
	}

	runner := grs.ExecRunner{}
	repos := defaultRepos(args)
	script := defaultScript(args)
	s := status.NewStatusboard(repos...)
	for _, repo := range s.Repos() {
		fmt.Printf("repos [%v] status is %v\n", repo, script(grs.Repo{Path:repo}, runner))
	}
}



// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepos(args Args) []grs.Repo {
	if len(args.repos) != 0 {
		grs.Debug("Using repos from CLI: %v", args.repos)
		return mkrepos(args.repos)
	}

	val, ok := os.LookupEnv("GRS_REPO")
	if ok {
		grs.Debug("Using repos from $GRS_REPO: %v", val)
		return mkrepos(val)
	}

	grs.Debug("Using default repos %v", os.ExpandEnv("$HOME" + string(os.PathSeparator) + "grstest"))
	return mkrepos(os.ExpandEnv("$HOME/grstest"))
}

func mkrepos(s string) []grs.Repo {
	res := make([]grs.Repo, 0, 1)
	for _, elem := range strings.Split(s, string(os.PathListSeparator)) {
		res = append(res, grs.Repo{Path:elem})
	}
	return res
}

func defaultScript(args Args) script.Script {
	grs.Debug("Using hard-coded script `GetRepoStatus`")
	return script.GetRepoStatus
}