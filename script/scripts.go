package script

import (
	"errors"
	"fmt"
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"os"
	"strconv"
	"strings"
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
	command := runner.Command(git, "show-ref", "-q", "HEAD")
	if _, err := command.CombinedOutput(); err != nil {
		rstat.Dir = status.DIR_INVALID
		return
	}
	rstat.Dir = status.DIR_VALID
}

// Fetch runs `git fetch`.
func Fetch(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat, repo grs.Repo) {
	if rstat.Dir != status.DIR_VALID {
		return
	}

	dbRepo := ctx.DB().FindOrCreateRepo(repo.Path)
	now := time.Now().Unix()
	if dbRepo.FetchedSec > (now - int64(ctx.MinFetchSec)) {
		return
	}

	git := ctx.GetGitExec()

	command := runner.Command(git, "fetch")
	if out, err := command.CombinedOutput(); err != nil {
		// fetch may have failed for common reasons, such as not adding yourxk ssh key to the agent
		grs.Debug("git fetch failed: %v\n%v", err, string(out))
		return
	}
	grs.Debug("git fetch ok: %v", repo.Path)
	dbRepo.FetchedSec = now
}

// GetRepoStatus gives a summary of the repo's status. Sets a number of `rstat` properties.
func GetRepoStatus(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat) {
	if rstat.Dir != status.DIR_VALID {
		return
	}

	git := ctx.GetGitExec()
	var command grs.Command
	var out []byte
	var err error

	command = runner.Command(git, "rev-parse", "@{upstream}")
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("GetRepoStatus: no upstream detected", err, string(out))
		rstat.Branch = status.BRANCH_UNTRACKED
		return
	}

	command = runner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
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
	local  int
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
	command := runner.Command(git, "ls-files", "--exclude-standard", "-om")
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

	command = runner.Command(git, "diff-index", "HEAD")
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
