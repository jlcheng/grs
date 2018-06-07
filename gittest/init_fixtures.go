package gittest

import (
	"fmt"
	"os"
)

// TestID:it_test_1 Sets up a git repository for it_test_1, rooted at tmpdir.
func InitTest1(tctx TestContext, tmpdir string) (err error) {
	if err := os.Chdir(tmpdir); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("error %v", r)
			}
		}
	}()

	git := tctx.git
	tctx.Exec(git, "init")
	tctx.TouchAndCommit(".gitignore", "Commit_A")
	tctx.TouchAndCommit("b.txt", "Commit_B")
	tctx.Exec(git, "checkout", "-b", "branch_A", "master~1")
	tctx.TouchAndCommit("c.txt", "Commit_C")
	tctx.Exec(git, "checkout", "master")
	tctx.TouchAndCommit("d.txt", "Commit_D")
	tctx.Exec(git, "checkout", "branch_A")
	tctx.TouchAndCommit("e.txt", "Commit_E")
	tctx.Exec(git, "checkout", "-b", "branch_B", "master~1")
	tctx.TouchAndCommit("f.txt", "Commit_F")
	tctx.TouchAndCommit("g.txt", "Commit_G")
	tctx.Exec(git, "checkout", "master")
	tctx.Exec(git, "merge", "branch_B", "-m", "merge branch_B onto master")
	tctx.Exec(git, "checkout", "branch_A")
	return err
}
