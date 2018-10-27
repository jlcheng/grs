// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string  // viper config
	verbose bool    // enables more logging
	daemon  bool    // runs in daemon mode
	refresh int     // how often to check for changes, in seconds
	forceMerge bool // ignore access time check when auto-merging
	repo string     // a single repo to process
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grs",
	Short: "grs performs two-way sync of Git repos",
	Long: "grs performs two-way sync of Git repos",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, _ []string) {
		var args = Args {
			verbose: verbose,
			daemon: daemon,
			refresh: viper.GetInt("refresh"),
			forceMerge: forceMerge,
			repos: viper.GetStringSlice("repos"),
		}
		if repo != "" {
			args.repos = []string{repo}
		}
		RunCli(args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grs.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbosity")
	rootCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "daemon mode")
	rootCmd.PersistentFlags().IntVarP(&refresh, "refresh", "t", 600, "how often to check for changes, in seconds")
	rootCmd.PersistentFlags().BoolVarP(&forceMerge, "force", "m", false, "ignore access time check when auto-merging")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "the repository to process")

	viper.BindPFlag("refresh", rootCmd.PersistentFlags().Lookup("refresh"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
