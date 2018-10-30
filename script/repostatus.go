package script

import (
	"errors"
	"fmt"
	"jcheng/grs/shexec"
	"strconv"
	"strings"
)

// GetRepoStatus() updates the status of a repository
func (s *Script) GetRepoStatus() {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil || repo.Dir != DIR_VALID {
		return
	}

	git := ctx.GetGitExec()
	var command shexec.Command
	var out []byte
	var err error

	command = ctx.CommandRunner.Command(git, "rev-parse", "@{upstream}")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("GetRepoStatus: no upstream detected. %s, %s", err, string(out))
		repo.Branch = BRANCH_UNTRACKED
		return
	}

	command = ctx.CommandRunner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("rev-list failed: %v\n%v", err, string(out))
		repo.Dir = DIR_INVALID
		return
	}
	diff, err := parseRevList(out)
	if err != nil {
		shexec.Debug("cannot parse `git rev-list...` output: %q", string(out))
		repo.Dir = DIR_INVALID
		return
	}

	shexec.Debug("CMD: git rev-list --left-right --count @{upstream}...HEAD")
	if diff.remote == 0 && diff.local == 0 {
		repo.Branch = BRANCH_UPTODATE
		return
	}
	if diff.remote > 0 && diff.local == 0 {
		repo.Branch = BRANCH_BEHIND
		return
	}
	if diff.remote == 0 && diff.local > 0 {
		repo.Branch = BRANCH_AHEAD
		return
	}
	if diff.remote > 0 && diff.local > 0 {
		repo.Branch = BRANCH_DIVERGED
		return
	}
	return
}

type remoteDiff struct {
	local  int
	remote int
}

func parseRevList(out []byte) (diff remoteDiff, err error) {
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