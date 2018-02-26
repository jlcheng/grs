package script

import (
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"fmt"
	"os"
	"strings"
	"errors"
	"strconv"
	"jcheng/grs/grsdb"
	"time"
)

// BeforeScript chdir to the repo directory and validates the repo. Sets rstat.Dir to `DIR_VALID` on success.
func BeforeScript(ctx *grs.AppContext, repo grs.Repo, runner grs.CommandRunner, rstat *status.RStat) {
	git := ctx.GetGitExec()

	if f, err := os.Stat(repo.Path); err != nil || !f.IsDir() {
		rstat.Dir = status.DIR_INVALID
		return
	}
	if err := os.Chdir(repo.Path); err != nil {
		rstat.Dir = status.DIR_INVALID
		return
	}
	command := *runner.Command(git, "show-ref", "-q", "HEAD")
	if _, err := command.CombinedOutput(); err != nil {
		rstat.Dir = status.DIR_INVALID
		return
	}
	rstat.Dir = status.DIR_VALID
}

// Fetch runs `git fetch`. Sets rstat.Dir to `DIR_INVALID` on error
func Fetch(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat, repo grs.Repo) {
	if rstat.Dir != status.DIR_VALID {
		return
	}

	var rIdx = -1
	db := ctx.DB()

	for i, r := range db.Repos {
		if r.Id == repo.Path {
			rIdx = i
			break
		}
	}
	if rIdx == -1 {
		tmp := grsdb.Repo{Id: repo.Path}
		db.Repos = append(db.Repos, tmp)
		rIdx = len(db.Repos) - 1
	}
	now := time.Now().Unix()
	if db.Repos[rIdx].FetchedSec > (now - int64(ctx.MinFetchSec)) {
		return
	}

	git := ctx.GetGitExec()

	command := *runner.Command(git, "fetch")
	if out, err := command.CombinedOutput(); err != nil {
		grs.Debug("fetch failed: %v\n%v", err, string(out))
		rstat.Dir = status.DIR_INVALID
		return
	}
	grs.Debug("git fetch for %v", repo.Path)

	db.Repos[rIdx].FetchedSec = time.Now().Unix()
}

// GetRepoStatus gives a summary of the repo's status. Sets a number of `rstat` properties.
func GetRepoStatus(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat) {
	if rstat.Dir != status.DIR_VALID {
		return
	}

	git := ctx.GetGitExec()

	rstat.Dir = status.DIR_VALID
	command := *runner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("rev-list failed: %v\n%v", err, string(out))
		rstat.Dir = status.DIR_INVALID
		return
	}
	diff, err := parseRevList(out)
	if err != nil {
		grs.Info("cannot parse `git rev-list...` output: %q", string(out))
		rstat.Dir = status.DIR_INVALID
		return
	}

	grs.Debug("CMD: git rev-list --left-right --count @{upstream}...HEAD")

	if diff.remote == 0 && diff.local == 0 {
		rstat.Branch = status.BRANCH_UPTODATE
		return
	}
	if diff.remote > 0 && diff.local == 0 {
		rstat.Branch = status.BRANCH_BEHIND
		return
	}
	if diff.remote == 0 && diff.local > 0 {
		rstat.Branch = status.BRANCH_AHEAD
		return
	}
	if diff.remote > 0 && diff.local > 0 {
		rstat.Branch = status.BRANCH_DIVERGED
		return
	}
	return
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

// GetIndexStatus sets the rstat.index property to modified if there are uncommited changes or if the index has been
// modified
func GetIndexStatus(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat) {
	if rstat.Dir != status.DIR_VALID {
		return
	}

	git := ctx.GetGitExec()

	rstat.Index = status.INDEX_UNKNOWN
	command := *runner.Command(git, "ls-files", "--exclude-standard", "-om")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("ls-files failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		rstat.Index = status.INDEX_MODIFIED
		return
	}

	command = *runner.Command(git, "diff-index", "HEAD")
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("diff-index failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		rstat.Index = status.INDEX_MODIFIED
		return
	}

	rstat.Index = status.INDEX_UNMODIFIED
}


type Script func(grs.Repo, grs.CommandRunner) status.RStat
