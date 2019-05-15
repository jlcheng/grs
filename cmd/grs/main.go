package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"jcheng/grs/ui"
	"os"
)

func main() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "grs",
		Short: "grs performs two-way sync of Git repos",
		Long:  "grs performs two-way sync of Git repos",
		// Uncomment the following line if your bare application
		// has an action associated with it:
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
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
