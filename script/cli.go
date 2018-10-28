package script

import (
	"fmt"
	"github.com/spf13/viper"
	"jcheng/grs/shexec"
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
	repoCfg    map[string]interface{}
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
		repoCfg:    viper.GetStringMap("repo_config"),
	}
	return args
}

func RunCli(args Args) {
	if args.verbose {
		shexec.SetLogLevel(shexec.DEBUG)
	}

	ctx := shexec.NewAppContextWithRunner(&shexec.ExecRunner{})
	repos := ReposFromStringSlice(args.repos, args.repoCfg)

	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	gui := NewGUI(args.daemon)
	syncController := NewSyncController(repos, ctx, gui)

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


// TODO: JCHENG unit test improvements
func ReposFromStringSlice(repos []string, repoCfg map[string]interface{}) []Repo {
	r := make([]Repo, len(repos))
	for idx, repoPath := range repos {
		r[idx] = Repo{Path: repoPath}
		repo := &r[idx]
		cfg, ok := GetStringMap(repoCfg, repoPath)
		if !ok {
			continue
		}
		if value, ok := GetBool(cfg, "push_allowed"); ok {
			repo.PushAllowed = value
		}
	}
	return r
}

func GetBool(stringMap map[string]interface{}, key string) (bool, bool) {
	value, ok := stringMap[key]
	if !ok {
		return false, false
	}
	boolv, ok := value.(bool)
	if !ok {
		return false, false
	}
	return boolv, true
}

func GetStringMap(stringMap map[string]interface{}, key string) (map[string]interface{}, bool) {
	value, ok := stringMap[key]
	if !ok {
		return nil, false
	}
	mapv, ok := value.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return mapv, true
}