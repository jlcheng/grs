package test

import (
	"testing"
	"jcheng/grs/config"
)

func TestConfigFromEnv(t *testing.T) {
	p := config.NewConfigParams()
	p.Env = "data/config.json"
	c, e := config.ReadConfig(p)
	verifyConfigJson(t, c, e)
}

func TestConfigFromUserConf(t *testing.T) {
	p := config.NewConfigParams()
	p.User = "data/config.json"
	c, e := config.ReadConfig(p)
	verifyConfigJson(t, c, e)
}

func TestPriorityEnv(t *testing.T) {
	p := config.NewConfigParams()
	p.Env = "data/config.json"
	p.User = "data/empty_config.json"
	c, e := config.ReadConfig(p)
	verifyConfigJson(t, c, e)
}

func verifyConfigJson(t *testing.T, c *config.Config, e error) {
	if e != nil {
		t.Error("error reading conf file", e)
		return
	}
	if c.Git != "/path/to/git" {
		t.Error("unexpected config.git", c.Git)
		return
	}
	if len(c.Repos) != 2 {
		t.Error("unexpected config.repos length", len(c.Repos))
		return
	}
	if c.Repos[0].Path != "rel/repo1" {
		t.Error("unexpected config.repos[0]", c.Repos[0])
		return
	}
	if c.Repos[1].Path != "/abs/repo2" {
		t.Error("unexpected config.repos[1]", c.Repos[1])
		return
	}
}