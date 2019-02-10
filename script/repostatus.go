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

	git := ctx.GitExec
	var command shexec.Command
	var out []byte
	var err error

	command = ctx.CommandRunner.Command(git, "rev-parse", "@{upstream}")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("GetRepoStatus: no upstream detected. %s, %s", err, strings.TrimSpace(string(out)))
		repo.Branch = BRANCH_UNTRACKED
		return
	}

	shexec.Debug("CMD: git rev-list --left-right --count @{upstream}...HEAD")
	command = ctx.CommandRunner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("git rev-list failed: %v\n%v", err, string(out))
		repo.Dir = DIR_INVALID
		return
	}
	shexec.Debug(strings.TrimSpace(string(out)))
	diff, err := parseRevList(out)
	if err != nil {
		shexec.Debug("cannot parse `git rev-list...` output: %q", string(out))
		repo.Dir = DIR_INVALID
		return
	}

	if diff.remote == 0 {
		if diff.local == 0 {
			repo.Branch = BRANCH_UPTODATE
		} else {
			repo.Branch = BRANCH_AHEAD
		}
	} else {
		if diff.local == 0 {
			repo.Branch = BRANCH_BEHIND
		} else {
			repo.Branch = BRANCH_DIVERGED
		}
	}
}

func (s *Script) GetCommitTime() {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil || repo.Dir != DIR_VALID {
		return
	}

	git := ctx.GitExec
	var command shexec.Command
	var out []byte
	var err error

	command = ctx.CommandRunner.Command(git, "log", "-1", "--format=%cr")
	shexec.Debug("CMD: git log -1 --format=%%cr")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("failed: %v\n%v\n", err, string(out))
		repo.CommitTime = "Unknown"
	}
	shexec.Debug(strings.TrimSpace(string(out)))
	repo.CommitTime = strings.Trim(string(out), "\n")
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
