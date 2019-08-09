package base

import (
	"github.com/spf13/viper"
)

func Version() string {
	version := viper.GetString("version")
	if version == "" {
		version = "Vx.x.x_12345"
	}
	return version
}

func SetVersion(version string) {
	viper.Set("version", version)
}
