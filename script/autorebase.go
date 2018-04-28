package script

import (
	"fmt"
	"jcheng/grs/config"
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"os"
	"path/filepath"
	"strings"
)

func AutoRebase(ctx *grs.AppContext, runner grs.CommandRunner, rstat *status.RStat, repo grs.Repo) {
	git := ctx.GetGitExec()
	r := *runner.Command(git, "")
	_ = r
	// Set up a working directory and update some sort of metadata object (grsdb)
	// for any repo that requires rebasing (branch == diverged):
	//  1. Set up a clone directory
	clnpath := ToClonePath(repo.Path)
	fmt.Printf("git clone %v %v\n", repo.Path, filepath.Dir(clnpath))

	//  2. Identify merge-base
	//  3. Identify the graph of child commits from merge-base to HEAD
	//  4. Abort if a node with multiple parents is identified; otherwise call it lineage
	//  5. Rebase current branch against each child in lineage
	//  6. Stop when conflict is detected
	//  7. Now the cloned branch is "as caught up as possible" against ${upstream}

}

func ToClonePath(repoPath string) string {
	return filepath.Join(config.UserPrefDir,
		"clones",
		strings.Replace(repoPath, string(os.PathSeparator), "_", -1))
}
