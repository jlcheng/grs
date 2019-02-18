package ui

import (
	"github.com/spf13/viper"
	"jcheng/grs/base"
	"jcheng/grs/script"
	"log"
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
	useCui     bool
}

// CliParse uses spf13/viper to create the program parameters
func CliParse() Args {
	// allow one to override the repo_config setting to run one-off tests
	repos := viper.GetStringSlice("repos")
	if repo := viper.GetString("repo"); repo != "" {
		repos = []string{repo}
	}

	var args = Args{
		verbose:    viper.GetBool("verbose"),
		daemon:     viper.GetBool("daemon"),
		refresh:    viper.GetInt("refresh"),
		forceMerge: viper.GetBool("merge-ignore-atime"),
		repos:      repos,
		repoCfgMap: parseRepoConfigMap(viper.Get("repo_config")),
		useCui:     viper.GetBool("use-cui"),
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
		base.SetLogLevel(base.DEBUG)
	}

	ctx := script.NewAppContext()
	repos := ReposFromStringSlice(args.repos, args.repoCfgMap)

	if len(repos) == 0 {
		log.Fatal("repos not specified")
	}


	var cliUI CliUI
	var err error
	if args.useCui {
		cliUI, err = NewConsoleUI()
	} else {
		cliUI, err = NewPrintUI()
	}
	if err != nil {
			log.Fatal("cannot initialize the terminal", err)
	}
	defer cliUI.Close()
	syncController := NewSyncController(repos, ctx, cliUI)
	syncController.Duration = time.Duration(args.refresh) * time.Second
	syncController.CliUIImpl()
}

// TODO: JCHENG unit test improvements
func ReposFromStringSlice(repos []string, repoCfg map[string]RepoConfig) []script.Repo {
	r := make([]script.Repo, len(repos))
	for idx, repoPath := range repos {
		r[idx] = script.Repo{Path: repoPath}
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
