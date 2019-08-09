package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"jcheng/grs/base"
	"jcheng/grs/ui"
	"os"
)

var Version string = "PLACE_HOLDER_VERSION"
var configFile = ""

func main() {
	base.SetVersion(Version)

	var rootCmd = &cobra.Command{
		Use:   "grs [-t 60] [--use-tui]",
		Short: "grs performs two-way sync of Git repos",
		Long: fmt.Sprintf(`grs %v
Grs performs two-way sync of Git repos
`, Version),
		Run: func(cmd *cobra.Command, _ []string) {
			ui.RunCli(ui.CliParse())
		},
	}

	cobra.OnInitialize(initConfig)
	pflag.StringVar(&configFile, "config", "", "The config file (default is $HOME/.grs.toml)")
	pflag.BoolP("verbose", "v", false, "Output verbose logs")
	pflag.IntP("refresh", "t", 600, "How often to poll for changes, in seconds")
	pflag.BoolP("merge-ignore-atime", "m", false, "Ignore access time check when auto-merging")
	pflag.StringP("repo", "r", "", "The repository to process")
	pflag.Bool("simple-ui", false, "Use a simple UI that does not put terminal into raw mode")
	pflag.String("log-file", "", "Where to write log messages. If not set, emit logs to stdout.")
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grs")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
