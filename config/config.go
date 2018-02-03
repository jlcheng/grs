package config

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Git string
	Repos []RepoConf
}

type RepoConf struct {
	Path string
}

type ConfigParams struct {
	Env string
	User string
}

func NewConfigParams() *ConfigParams {
	return &ConfigParams{
		Env: os.Getenv("GRS_CONF"),
		User: filepath.Join(os.ExpandEnv("${HOME}"),".grs.json"),
	}
}

func GetCurrConfig(p *ConfigParams) (*Config, error) {
	if len(p.Env) != 0 {
		return readConfFile(p.Env)
	}
	if _, err := os.Stat(p.User); err == nil {
		return readConfFile(p.User)
	}
	return &Config{}, nil
}

func readConfFile(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		return nil, err
	}
	var c = &Config{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
