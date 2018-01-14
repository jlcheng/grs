package script

import (
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"fmt"
	"os"
	"strings"
	"errors"
	"strconv"
)

func GetRepoStatus(repo grs.Repo, runner grs.CommandRunner) status.RepoStatus {
	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		return status.INVALID
	}
	if err := os.Chdir(repo.Path); err != nil {
		return status.INVALID
	}
	command := *runner.Command("git", "rev-list", "--left-right", "--count", "@{u}..@")
	var out []byte;
	var err error;
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("rev-list failed: %v\n%v\n", err, string(out))
		return status.UNKNOWN
	}
	diff, err := parseRevList(out)
	if err != nil {
		grs.Debug("cannot parse `git rev-list...` output: %q", string(out))
		return status.UNKNOWN
	}
	if diff.remote == 0 && diff.local == 0 {
		return status.LATEST
	}
	if diff.remote > 0 && diff.local == 0 {
		return status.BEHIND
	}
	if diff.remote == 0 && diff.local > 0 {
		return status.AHEAD
	}
	if diff.remote > 0 && diff.local > 0 {
		return status.DIVERGED
	}
	return status.UNKNOWN
}

func GetWorkingDirStatus(repo grs.Repo, runner grs.CommandRunner) status.RepoStatus {
	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		return status.INVALID
	}
	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		return status.INVALID
	}
	if err := os.Chdir(repo.Path); err != nil {
		return status.INVALID
	}
	command := *runner.Command("git", "rev-list", "--left-right", "--count", "@\\{u\\}..@")
	var err error;
	if _, err = command.CombinedOutput(); err != nil {
		grs.Debug("command failed: %v\n", err)
		return status.UNKNOWN
	}

	return status.UNKNOWN
}


type RemoteDiff struct {
	local int
	remote int
}

func parseRevList(out []byte) (diff RemoteDiff, err error) {
	str := strings.TrimSpace(string(out))
	tokens := strings.Split(str, "\t")
	if len(tokens) != 2 {
		return diff, errors.New(fmt.Sprintf("expected token count=2, got [%v]", len(tokens)))
	}
	diff.remote, err = strconv.Atoi(tokens[0])
	if err != nil {
		return diff, err
	}
	diff.local, err = strconv.Atoi(tokens[1])
	if err != nil {
		return diff, err
	}
	return diff, nil
}

type Script func(grs.Repo, grs.CommandRunner) status.RepoStatus