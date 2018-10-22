package grs

import (
	"bytes"
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

func ReposFromString(input string) []status.Repo {
	tokens := strings.Split(input, string(os.PathListSeparator))
	r := make([]status.Repo, len(tokens))
	for idx, elem := range tokens {
		r[idx] = status.Repo{Path: elem}
	}
	return r
}

func ReposFromStringSlice(input []string) []status.Repo {
	r := make([]status.Repo, len(input))
	for idx, elem := range input {
		r[idx] = status.Repo{Path: elem}
	}
	return r
}
