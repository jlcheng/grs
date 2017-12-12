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
	if len(args.repo) == 0 {
		args.repo = os.ExpandEnv("$HOME/github/test")
	}
	if args.verbose {
		fmt.Printf("using repo: %v\n", args.repo)
	}

	repo := gitscripts.Repo{args.repo}
	c, err := gitscripts.Status(repo)
	if err != nil {
		fmt.Println("err ", err)
		return
	}

	fmt.Println(c.String())
}


