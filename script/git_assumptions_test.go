// git_assumptions_test.go contains a set of 'tests' used to explore the behavior of Git commands. For GRS to work, these tests
// likely must 'pass'. This is a collection of assumptions I make about Git's behaviors.0
package script

import (
	"os"
	"strings"
	"testing"
)


// Justification for autorebase. Documents the behavior of `git pull --rebase` when there is a conflict.
// Here, origin and master each has one conflict-free change and one conflicting change. The desirable state is for
// master (local) to end up with:
// 1) Conflict-free change automatically resolved
// 2) Conflicting change at the tip of the commit log
// 3) Showing that the local master and origin master has diverged
//
// However, this experiment shows that git `git pull --rebase` will leave the index in a "must resolve conflict"
// mode. This behavior is the motivation for the "autorebase" functionality of Grs.
func TestRebasePullConflict(t *testing.T) {
	// === START: repo initialization ===
	// TODO JCHENG take this initialization section and move it into GitHelper
	// Creates a bare repo named '$tmp_dir/source' and a working directory named `$tmp_dir/dest`
	const TEST_LABEL = "TestRebasePullConflict"
	oldwd, tmpdir := MkTmpDir(t, TEST_LABEL, TEST_LABEL)
	defer CleanTmpDir(t, oldwd, tmpdir, TEST_LABEL)
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

	exec := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	git := exec.Git()
	exec.Mkdir("source")
	exec.Chdir("source")
	exec.Exec(git, "init", "--bare")
	exec.Chdir("..")
	exec.Exec(git, "clone", "source", "dest")

	exec.Chdir("dest")
	exec.TouchAndCommit(".gitignore", "Commit_A")

	if exec.Err() != nil {
		t.Fatal("test setup failed\n" + exec.ErrString())
	}
	// === END: repo initialization ===

	exec.SetContents("a_1.txt", "1\n2\n3\n")
	exec.Add("a_1.txt")
    exec.Commit("conflict-free change on a")
	exec.Exec(git, "push", "origin")

	exec.Exec(git, "reset", "--hard", "HEAD~1")
    exec.TouchAndCommit("be_1.txt", "conflict-free change on b")
	exec.SetContents("a_1.txt", "1\n3\n")
	exec.Add("a_1.txt")
	exec.Commit("conflict on b")
	exec.Exec(git, "pull", "--rebase", "-v")

	conflict := strings.Contains(exec.ErrString(), ": Merge conflict in a_1.txt")
	if !conflict {
		t.Fatal("Expected conflict, got the following instead.",  "\n\n"+exec.ErrString())
	}

}


