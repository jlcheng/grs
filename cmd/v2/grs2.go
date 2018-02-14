package main

import (
	"jcheng/grs/grs"
	"os"
	"flag"
	"fmt"
	"jcheng/grs/script"
	"strings"
	"jcheng/grs/status"
	"jcheng/grs/config"
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
	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	status := status.NewStatusboard(repos...)
	for _, elem := range status.Repos() {
		fmt.Printf("repo [%v]\n", elem)
		repo := grs.Repo{Path:elem}
		rstat := script.Fetch(repo, runner)
		if rstat.Dir == 0 {
			rstat = script.GetRepoStatus(repo, runner)
			fmt.Printf("repos [%v] status is %v\n", repo, rstat.Branch)
		}
	}
}



// defaultRepo returns a Repo based on CLI args, Env variable, then defaults to "$HOME/grstest"
func defaultRepos(args Args) []grs.Repo {
	if len(args.repos) != 0 {
		grs.Debug("Using repos from CLI: %v", args.repos)
		return mkrepos(args.repos)
	}

	p := config.NewConfigParams()
	if c, _ := config.GetCurrConfig(p); c != nil {
		return grs.ReposFromConf(c.Repos)
	}
	return []grs.Repo{}
}

func mkrepos(s string) []grs.Repo {
	res := make([]grs.Repo, 0, 1)
	for _, elem := range strings.Split(s, string(os.PathListSeparator)) {
		res = append(res, grs.Repo{Path:elem})
	}
	return res
}
