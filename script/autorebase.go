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
)

const (
	CLONE_BASEDIR = "clones"
)

func AutoRebase(ctx *grs.AppContext, repo grs.Repo, runner grs.CommandRunner, rstat *status.RStat, clone bool) {
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
				fmt.Println(err)
				return
			}
		}

		err = grsio.CopyDir(repo.Path, clnpath)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Chdir(clnpath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//  2. Identify merge-base
	git := ctx.GetGitExec()
	var cmd = runner.Command(git, "merge-base", "HEAD", "master")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}
	mergeBase := strings.TrimSpace(string(bytes))

	//  3. Identify the graph of child commits from merge-base to HEAD
	cmd = runner.Command(git, "rev-list", "master", "^"+mergeBase)
	bytes, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}

	revlist := strings.Split(strings.TrimSpace(string(bytes)),"\n")
	//  5. Rebase current branch against each child in lineage
	for i := len(revlist)-1; i >=0; i-- {
		commit := revlist[i]
		cmd = runner.Command(git, "rebase", commit)
		bytes, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//  6. Stop when conflict is detected



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