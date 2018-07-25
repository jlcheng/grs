package grs

import (
	"bytes"
	"jcheng/grs/config"
	"jcheng/grs/status"
	"os"
	"os/exec"
	"strings"
)

type Result struct {
	delegate *exec.Cmd
	Stdout   string
}

func (cmd *Result) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}

func ReposFromConf(rc []config.RepoConf) []status.Repo {
	var r = make([]status.Repo, len(rc))
	for idx, elem := range rc {
		r[idx] = status.Repo{Path: elem.Path}
	}
	return r
}

func ReposFromString(input string) []status.Repo {
	tokens := strings.Split(input, string(os.PathListSeparator))
	r := make([]status.Repo, len(tokens))
	for idx, elem := range tokens {
		r[idx] = status.Repo{Path: elem}
	}
	return r
}
