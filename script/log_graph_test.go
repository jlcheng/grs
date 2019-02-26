package script

import (
	"reflect"
	"testing"
)

// Verifies that Scirpt.LogGraph produces a graph object
func TestLogGraph(t *testing.T) {
	const TEST_ID = "TestLogGraph"
	tmpdir, cleanup := MkTmpDir1(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	repo := NewRepo(gh.Getwd())
	repo.PushAllowed = true
	s := NewScript(NewAppContext(), repo)
	s.BeforeScript()

	/* Creates commit log with this form
    a--b---c---e
     \  \     /
      \  d---e
       \
        g---h
	*/
	gh.TouchAndCommit("A.txt", "Commit_A")
	gh.GitExec("tag", "Tag_A")
	gh.TouchAndCommit("B.txt", "Commit_B")
	gh.GitExec("tag", "Tag_B")
	gh.TouchAndCommit("C.txt", "Commit_C")
	gh.GitExec("checkout", "-b", "source_2", "Tag_B")
	gh.TouchAndCommit("D.txt", "Commit_D")
	gh.TouchAndCommit("E.txt", "Commit_E")
	gh.GitExec("checkout", "master")
	gh.GitExec("merge", "source_2", "-m", "Merge_C_E")

	expected := map[string][]string{
		"Merge_C_E": {"Commit_C", "Commit_E"},
		"Commit_E": {"Commit_D"},
		"Commit_D": {"Commit_B"},
		"Commit_C": {"Commit_B"},
		"Commit_B": {"Commit_A"},
		"Commit_A": {"Initial commit"},
		"Initial commit": {},
	}
	if lg, err := s.LogGraph(); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(expected, lg) {
		t.Fatal("did not match expected commit grpah")
	}

	gh.GitExec("push")
	gh.GitExec("reset", "--hard", "Tag_A")
	gh.TouchAndCommit("G.txt", "Commit_G")
	gh.TouchAndCommit("H.txt", "Commit_H")
}

