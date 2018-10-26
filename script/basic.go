package script

import (
	"bytes"
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

func ReposFromString(input string) []Repo {
	tokens := strings.Split(input, string(os.PathListSeparator))
	r := make([]Repo, len(tokens))
	for idx, elem := range tokens {
		r[idx] = Repo{Path: elem}
	}
	return r
}

func ReposFromStringSlice(input []string) []Repo {
	r := make([]Repo, len(input))
	for idx, elem := range input {
		r[idx] = Repo{Path: elem}
	}
	return r
}
