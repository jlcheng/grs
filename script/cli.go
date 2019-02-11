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
	concurrent bool
}

func CliParse(verbose bool, daemon bool, refresh int, forceMerge bool, repo string, concurrent bool) Args {
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
		concurrent: concurrent,
	}
	return args
}

func parseRepoConfigMap(obj interface{}) map[string]RepoConfig {
	if sliceIfc, ok := obj.([]interface{}); ok {
		sliceStringMap := ToSliceStringMap(sliceIfc)
		return ToRepoConfigMap(sliceStringMap)
	}
	return make(map[string]RepoConfig)
}

type RepoConfig struct {
	pushAllowed bool
}

func RunCli(args Args) {
	if args.verbose {
		shexec.SetLogLevel(shexec.DEBUG)
	}

	ctx := shexec.NewAppContext()
	repos := ReposFromStringSlice(args.repos, args.repoCfgMap)

	if len(repos) == 0 {
		fmt.Println("repos not specified")
		os.Exit(1)
	}

	gui := NewGUI(args.daemon)
	syncController := NewSyncController(repos, ctx, gui, args.concurrent)

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

func GetBool(stringMap map[string]interface{}, key string, fallback bool) bool {
	value, ok := stringMap[key]
	if !ok {
		return fallback
	}
	boolv, ok := value.(bool)
	if !ok {
		return fallback
	}
	return boolv
}

func GetString(stringMap map[string]interface{}, key string, fallback string) string {
	value, ok := stringMap[key]
	if !ok {
		return fallback
	}
	stringv, ok := value.(string)
	if !ok {
		return fallback
	}
	return stringv
}

// Asserts that the given value is a slice of []map[string]interface{}, raising an error if not
func ToSliceStringMap(input []interface{}) []map[string]interface{} {
	emptySlice := make([]map[string]interface{}, 0)
	var output = make([]map[string]interface{}, len(input))
	for i := 0; i < len(output); i++ {
		elem, ok := input[i].(map[string]interface{})
		if !ok {
			return emptySlice
		}
		output[i] = elem
	}
	return output
}

func ToRepoConfigMap(input []map[string]interface{}) map[string]RepoConfig {
	var output = make(map[string]RepoConfig)
	for i := 0; i < len(input); i++ {
		rawMap := input[i]
		if repoID := GetString(rawMap, "id", ""); repoID != "" {
			var repoConfig = RepoConfig{}
			repoConfig.pushAllowed = GetBool(rawMap, "push_allowed", false)
			output[repoID] = repoConfig
		}
	}
	return output
}
