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

	gui := NewGUI(args.daemon)
	syncController := NewSyncController(repos, ctx, gui)

	if args.useCui {
	}

	if args.useCui {
		cui := NewCuiGUI()
		// Sets up a `done` channel. The `done` channel can only be closed by the terminal UI. Closing the channel
		// will cascade to various goroutines that coordinate the UI and application code.
		if err := cui.Init(); err != nil {
			log.Fatal("cannot initialize the terminal", err)
		}
		defer cui.Close()
		// SyncController needs to know about the terminal UI to tell it to redraw
		syncController.Cui = cui
		// SyncController needs to know about the refresh interval to know often it should run the `sync repo` code
		syncController.Duration = time.Duration(args.refresh) * time.Second
		// Gets reference of the `done` channel so the main goroutine can block until the UI receives a `done` signal
		done := cui.done

		// starts two goroutines
		//  syncer: a tick event triggers `sync repo` code, which outputs a `repos synced` event
		//  uiForwarder: a `repos synced` event triggers code to use the cui API to redraw the terminal UI
		go syncController.RunLoops()

		// starts goroutine for terminal UI
		go cui.MainLoop()

		// blocks until the UI has been stopped
		<- done
	} else {
		// run at least once
		syncController.Run()
		if args.daemon {
			ticker := time.NewTicker(time.Duration(args.refresh) * time.Second)
			defer ticker.Stop()
			if args.useCui {
				DAEMON_LOOP:
				for !syncController.Cui.stopped {
					select {
					case <-ticker.C:
						syncController.Run()
					case <-syncController.Cui.done:
						break DAEMON_LOOP
					}
				}

			} else {
				for {
					select {
					case <-ticker.C:
						syncController.Run()
					}
				}
			}
		}
	}

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
