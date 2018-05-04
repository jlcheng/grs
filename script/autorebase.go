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
)

func AutoRebase(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat, repo grs.Repo) {
	// Set up a working directory and update some sort of metadata object (grsdb)
	// for any repo that requires rebasing (branch == diverged):
	//  1. Set up a clone directory
	clnpath := ToClonePath(repo.Path)
	_, err := os.Stat(clnpath)
	if os.IsNotExist(err) {
		if err := CreateCloneDir(); err != nil {
			fmt.Println(err)
			return
		}

		err := grsio.CopyDir(repo.Path, clnpath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//  2. Identify merge-base
	//  3. Identify the graph of child commits from merge-base to HEAD
	//  4. Abort if a node with multiple parents is identified; otherwise call it lineage
	//  5. Rebase current branch against each child in lineage
	//  6. Stop when conflict is detected
	//  7. Now the cloned branch is "as caught up as possible" against ${upstream}

}

func ToClonePath(repoPath string) string {
	h := sha1.New()
	io.WriteString(h, repoPath)
	d := fmt.Sprintf("%x", h.Sum(nil))
	return filepath.Join(config.UserPrefDir, "clones", d)
}

func CreateCloneDir() error {
	d := filepath.Join(config.UserPrefDir, "clones")
	return os.Mkdir(d, 0755)
}
