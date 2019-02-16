package script

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Repo struct {
	Path        string
	Dir         Dirstat
	Branch      Branchstat
	Index       Indexstat
	CommitTime  string
	PushAllowed bool
}

func NewRepo(path string) *Repo {
	return &Repo{
		Path:       path,
		Dir:        DIR_INVALID,
		Branch:     BRANCH_UNKNOWN,
		Index:      INDEX_UNKNOWN,
		CommitTime: "",
	}
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
