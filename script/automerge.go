package script

import (
	"errors"
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"os"
	"path/filepath"
	"time"
)

// AutoFFMerge runs `git merge --ff-only...` when the branch is behind and unmodified
func (s *Script) AutoFFMerge() bool {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil ||
		repo.Dir != status.DIR_VALID ||
		repo.Branch != status.BRANCH_BEHIND ||
		repo.Index != status.INDEX_UNMODIFIED {
		return false
	}

	git := ctx.GetGitExec()

	command := ctx.CommandRunner.Command(git, "merge", "--ff-only", "@{upstream}")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("git merge failed: %v\n%v\n", err, string(out))
		return false
	}
	return true
}

var lastActivityFiles = []string{"HEAD", "COMMIT_EDITMSG", "ORIG_HEAD", "index", "config"}

// GetActivityTime gets the estimated "last modified time" of a repo
func GetActivityTime(repo string) (time.Time, error) {
	var atime time.Time
	if f, err := os.Stat(repo); err != nil || !f.IsDir() {
		return atime, errors.New("%v is not a directory")
	}
	for _, f := range lastActivityFiles {
		fn := filepath.Join(repo, ".git", f)
		if finfo, err := os.Stat(fn); err == nil {
			if finfo.ModTime().After(atime) {
				atime = finfo.ModTime()
			}
		}
	}
	return atime, nil
}
