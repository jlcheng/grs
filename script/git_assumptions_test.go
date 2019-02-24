// git_assumptions_test.go contains a set of 'tests' used to explore the behavior of Git commands. For GRS to work, these tests
// likely must 'pass'. This is a collection of assumptions I make about Git's behaviors.0
package script

import (
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

/*
Assume A, B, C are conflct-free changes, but D and E conflicts.

    A--B--D  origin
     \
      C--E   local

AutoRebase changes local to this. Git's rebase does not support this.

    A--B--D    origin
        \
         C--E  local

Not included here, but rebase using `git pull --rebase -s recursive -X ours` does not yield what I expect either:

    A--B--D    origin
        \
         D--C  local, E got lost

*/

func TestRebasePullConflict(t *testing.T) {
	tmpdir, cleanup := MkTmpDir1(t, "TestRebasePullConflict")
	defer cleanup()
	exec := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	exec.NewRepoPair(tmpdir)
	git := exec.Git()

	exec.Chdir("dest")
	exec.TouchAndCommit("B.txt", "B: conflict-free change on origin")
	exec.SetContents("conflict.txt", "1\n2\n3\n")
	exec.Add("conflict.txt")
	exec.Commit("D: conflicting change on origin")
	exec.Exec(git, "push", "origin")

	exec.Exec(git, "reset", "--hard", "HEAD~1")
    exec.TouchAndCommit("C.txt", "C: conflict-free change on local")
	exec.SetContents("conflict.txt", "1\n3\n")
	exec.Add("conflict.txt")
	exec.Commit("E: conflicting change on local")
	exec.Exec(git, "pull", "--rebase", "-v")

	conflict := strings.Contains(exec.ErrString(), ": Merge conflict in conflict.txt")
	if !conflict {
		t.Fatal("Expected conflict, got the following instead.",  "\n\n"+exec.ErrString())
	}
}


