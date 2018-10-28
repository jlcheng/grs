package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"jcheng/grs/shexec"
	"jcheng/grs/script"
	"os"
	"time"
)

type Args struct {
	repos      []string
	verbose    bool
	command    string
	daemon     bool
	refresh    int
	forceMerge bool
	repoConf   map[string]interface{}
}

func CliParse(verbose bool, daemon bool, refresh int, forceMerge bool, repo string) Args {
	// command line arg takes precedence over repos
	repos := viper.GetStringSlice("repos")
	if repo != "" {
		repos = []string{repo}
	}

	var args = Args{
		verbose:    verbose,
		daemon:     daemon,
		refresh:    viper.GetInt("refresh"),
		forceMerge: forceMerge,
		repos:      repos,
		repoConf:   viper.GetStringMap("repo_config"),
	}
	return args
}

func RunCli(args Args) {
	if args.verbose {
		shexec.SetLogLevel(shexec.DEBUG)
	}

	ctx := shexec.NewAppContextWithRunner(&shexec.ExecRunner{})
	repos := script.ReposFromStringSlice(args.repos)

	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	gui := script.NewGUI(args.daemon)
	syncController := script.NewSyncController(repos, ctx, gui)

	// run at least once
	syncController.Run()
	if args.daemon {
		ticker := time.NewTicker(time.Duration(args.refresh) * time.Second)
		defer ticker.Stop() // remove? not strictly necessary as we don't offer a way to gracefully shutdown

		// use Ctrl-C to stop this program
		for {
			select {
			case <-ticker.C:
				syncController.Run()
			}
		}
	}
}
