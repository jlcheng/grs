package config

import (
	"os"
	"path/filepath"
)

var UserPrefDir = filepath.Join(os.ExpandEnv("${HOME}"), ".grs.d")
var UserConf = filepath.Join(UserPrefDir, "config.json")
var UserDB = filepath.Join(UserPrefDir, "grs.db")
var UserDBName = "grs.db"

// SetupUser creates the $HOME/.grs.d directory if needed.
func SetupUserPrefDir(basedir string) error {
	if _, err := os.Stat(basedir); err != nil {
		if err = os.Mkdir(basedir, 700); err != nil {
			return err
		}
	}
	return nil
}
