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

func Fetch(repo grs.Repo, runner grs.CommandRunner) (rstat status.RStat) {
	ctx := grs.GetContext()
	git := ctx.GetGitExec()

	rstat = status.NewRStat()
	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	if err := os.Chdir(repo.Path); err != nil {
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	rstat.Dir = status.DIR_VALID
	command := *runner.Command(git, "fetch")
	var out []byte;
	var err error;
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("fetch failed: %v\n%v\n", err, string(out))
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	return rstat
}

func GetRepoStatus(repo grs.Repo, runner grs.CommandRunner) (rstat status.RStat) {
	ctx := grs.GetContext()
	git := ctx.GetGitExec()

	rstat = status.NewRStat()
	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	if err := os.Chdir(repo.Path); err != nil {
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	rstat.Dir = status.DIR_VALID
	command := *runner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}..HEAD")
	var out []byte;
	var err error;
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("rev-list failed: %v\n%v\n", err, string(out))
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	diff, err := parseRevList(out)
	if err != nil {
		grs.Debug("cannot parse `git rev-list...` output: %q", string(out))
		rstat.Dir = status.DIR_INVALID
		return rstat
	}
	if diff.remote == 0 && diff.local == 0 {
		rstat.Branch = status.BRANCH_UPTODATE
		return rstat
	}
	if diff.remote > 0 && diff.local == 0 {
		rstat.Branch = status.BRANCH_BEHIND
		return rstat
	}
	if diff.remote == 0 && diff.local > 0 {
		rstat.Branch = status.BRANCH_AHEAD
		return rstat
	}
	if diff.remote > 0 && diff.local > 0 {
		rstat.Branch = status.BRANCH_DIVERGED
		return rstat
	}
	return rstat
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

type Script func(grs.Repo, grs.CommandRunner) status.RStat