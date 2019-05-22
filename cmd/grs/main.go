package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"jcheng/grs/ui"
	"os"
)

var Version string = "PLACE_HOLDER_VERSION"

func main() {
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
	pflag.String("config", "", "config file (default is $HOME/.grs.toml)")
	pflag.BoolP("verbose", "v", false, "output verbose logs")
	pflag.IntP("refresh", "t", 600, "how often to poll for changes, in seconds")
	pflag.BoolP("merge-ignore-atime", "m", false, "ignore access time check when auto-merging")
	pflag.StringP("repo", "r", "", "the repository to process")
	pflag.Bool("use-tui", false, "use the experiment text-based UI")
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
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
