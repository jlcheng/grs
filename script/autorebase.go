package script

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"os"
	"path/filepath"
	"strings"
)

const (
	CLONE_BASEDIR = "clones"
)

func AutoRebase(ctx *grs.AppContext, runner grs.CommandRunner) error {
	//  2. Identify merge-base
	git := ctx.GetGitExec()
	p := "@{upstream}"
	cmd := runner.Command(git, "merge-base", "HEAD", p)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	mergeBase := strings.TrimSpace(string(bytes))

	//  3. Identify the graph of child commits from merge-base to HEAD
	cmd = runner.Command(git, "rev-list", p, "^"+mergeBase)
	bytes, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	revlist := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	//  5. Rebase current branch against each child in lineage
	for i := len(revlist) - 1; i >= 0; i-- {
		commit := revlist[i]
		cmd = runner.Command(git, "rebase", commit)
		_, err1 := cmd.CombinedOutput()
		if err1 != nil {
			cmd = runner.Command(git, "rebase", "--abort")
			bytes2, err2 := cmd.CombinedOutput()
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s", err2, string(bytes2)))
			}
		}
	}

	//  6. Stop when conflict is detected
	return nil
}

func ToClonePath(repoPath string) string {
	h := sha1.New()
	io.WriteString(h, repoPath)
	d := fmt.Sprintf("%x", h.Sum(nil))
	return filepath.Join(GetCloneBaseDir(), d)
}

func CreateCloneDir() error {
	return os.Mkdir(GetCloneBaseDir(), 0755)
}

func GetCloneBaseDir() string {
	return filepath.Join(config.UserPrefDir, CLONE_BASEDIR)
}
