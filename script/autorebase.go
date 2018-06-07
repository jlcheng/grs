package script

import (
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"path/filepath"
	"jcheng/grs/grsio"
	"fmt"
	"crypto/sha1"
	"io"
	"os"
	"strings"
	"errors"
)

const (
	CLONE_BASEDIR = "clones"
)

func AutoRebase(ctx *grs.AppContext, repo grs.Repo, runner grs.CommandRunner, rstat *status.RStat, clone bool) error {
	// Set up a working directory and update some sort of metadata object (grsdb)
	// for any repo that requires rebasing (branch == diverged):
	//  1. Set up a clone directory
	if (clone) {
		clnpath := ToClonePath(repo.Path)
		clndir, err := os.Stat(clnpath)
		if clndir != nil {
			os.RemoveAll(clndir.Name())
		}

		_, err = os.Stat(GetCloneBaseDir())
		if os.IsNotExist(err) {
			if err := CreateCloneDir(); err != nil {
				return err
			}
		}

		err = grsio.CopyDir(repo.Path, clnpath)
		if err != nil {
			return err
		}

		err = os.Chdir(clnpath)
		if err != nil {
			return err
		}
	}

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
	revlist := strings.Split(strings.TrimSpace(string(bytes)),"\n")

	//  5. Rebase current branch against each child in lineage
	for i := len(revlist)-1; i >=0; i-- {
		commit := revlist[i]
		cmd = runner.Command(git, "rebase", commit)
		bytes, err = cmd.CombinedOutput()
		if err != nil {
			cmd = runner.Command(git, "rebase", "--abort")
			bytes, err = cmd.CombinedOutput()
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s", err, string(bytes)))
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