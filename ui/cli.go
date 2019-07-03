package ui

import (
	"github.com/spf13/viper"
	"jcheng/grs/base"
	"jcheng/grs/script"
	"jcheng/grs/shexec"
	"log"
	"time"
)

type Args struct {
	forceMerge bool
	logFile    string
	refresh    int
	repoCfgMap map[string]RepoConfig
	repos      []string
	simpleUI   bool
	verbose    bool
}

// CliParse uses spf13/viper to create the program parameters
func CliParse() Args {
	// allow one to override the repo_config setting to run one-off tests
	repos := viper.GetStringSlice("repos")
	if repo := viper.GetString("repo"); repo != "" {
		repos = []string{repo}
	}

	var args = Args{
		forceMerge: viper.GetBool("merge-ignore-atime"),
		logFile:    viper.GetString("log-file"),
		refresh:    viper.GetInt("refresh"),
		repoCfgMap: parseRepoConfigMap(viper.Get("repo_config")),
		repos:      repos,
		simpleUI:   viper.GetBool("simple-ui"),
		verbose:    viper.GetBool("verbose"),
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
	if args.logFile != "" {
		if err := base.SetLogFile(args.logFile); err != nil {
			log.Fatal(err)
		}
	}

	grsRepos := InitGrsRepos(args.repos, args.repoCfgMap)

	cliUI := InitCliUI(args.simpleUI)
	defer cliUI.Close()

	refreshInterval := time.Duration(args.refresh) * time.Second
	syncController := NewSyncController(grsRepos, cliUI, refreshInterval)
	syncController.Run()
}

func InitGrsRepos(repos []string, repoCfgMap map[string]RepoConfig) []script.GrsRepo {
	grsRepos := GrsRepos(repos, repoCfgMap)

	if len(grsRepos) == 0 {
		log.Fatal("repos not specified")
	}
	return grsRepos
}

func InitCliUI(simpleUI bool) CliUI {
	var cliUI CliUI
	var err error
	if simpleUI {
		cliUI, err = NewPrintUI()
	} else {
		cliUI, err = NewConsoleUI()
	}
	if err != nil {
		log.Fatal("cannot initialize the terminal", err)
	}
	return cliUI
}

// GrsRepos parses a list of paths and a map of path configurations to derive a list of GrsRepo objects
func GrsRepos(paths []string, repoCfg map[string]RepoConfig) []script.GrsRepo {
	commandRunner := &shexec.ExecRunner{}
	repos := make([]script.GrsRepo, len(paths))
	for idx, path := range paths {
		config := repoCfg[path]
		pushAllowed := config.pushAllowed
		repos[idx] = script.NewGrsRepo(
			script.WithLocalGrsRepo(path),
			script.WithPushAllowed(pushAllowed),
			script.WithCommandRunnerGrsRepo(commandRunner),
		)
	}
	return repos
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
