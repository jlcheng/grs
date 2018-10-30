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
	repoCfgMap map[string]RepoConfig
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
		repoCfgMap: parseRepoConfigMap(viper.Get("repo_config")),
	}
	return args
}

func parseRepoConfigMap(obj interface{}) map[string]RepoConfig {
	var ok bool
	var sliceIfc []interface{}
	retval := make(map[string]RepoConfig)
	if sliceIfc, ok = obj.([]interface{}); !ok {
		return retval
	}

	for _, ifc := range sliceIfc {
		var ifcmap map[string]interface{}
		if ifcmap, ok = ifc.(map[string]interface{}); !ok {
			return retval
		}
		var repoPath string
		var pushAllowed bool
		if repoPath, ok = GetString(ifcmap, "id"); !ok {
			continue
		}
		if pushAllowed, ok = GetBool(ifcmap, "push_allowed"); !ok {
			continue
		}
		retval[repoPath] = RepoConfig{pushAllowed: pushAllowed}
	}
	return retval
}

type RepoConfig struct {
	pushAllowed bool
}

func RunCli(args Args) {
	if args.verbose {
		shexec.SetLogLevel(shexec.DEBUG)
	}

	ctx := shexec.NewAppContextWithRunner(&shexec.ExecRunner{})
	repos := ReposFromStringSlice(args.repos, args.repoCfgMap)

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
func ReposFromStringSlice(repos []string, repoCfg map[string]RepoConfig) []Repo {
	r := make([]Repo, len(repos))
	for idx, repoPath := range repos {
		r[idx] = Repo{Path: repoPath}
		repo := &r[idx]

		config, ok := repoCfg[repoPath]
		if !ok {
			continue
		}
		repo.PushAllowed = config.pushAllowed
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

func GetString(stringMap map[string]interface{}, key string) (string, bool) {
	value, ok := stringMap[key]
	if !ok {
		return "", false
	}
	stringv, ok := value.(string)
	if !ok {
		return "", false
	}
	return stringv, true
}
