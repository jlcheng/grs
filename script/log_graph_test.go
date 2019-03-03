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

	// Creates commit log with this form
	//  a--b---c---e
	//      \     /
	//       d---e
	gh.TouchAndCommit("A.txt", "Commit_A")
	gh.RunGit("tag", "Tag_A")
	gh.TouchAndCommit("B.txt", "Commit_B")
	gh.RunGit("tag", "Tag_B")
	gh.TouchAndCommit("C.txt", "Commit_C")
	gh.RunGit("checkout", "-b", "source_2", "Tag_B")
	gh.TouchAndCommit("D.txt", "Commit_D")
	gh.TouchAndCommit("E.txt", "Commit_E")
	gh.RunGit("checkout", "master")
	gh.RunGit("merge", "source_2", "-m", "Merge_C_E")

	expected := LogGraph(map[string][]string{
		"Merge_C_E": {"Commit_C", "Commit_E"},
		"Commit_E": {"Commit_D"},
		"Commit_D": {"Commit_B"},
		"Commit_C": {"Commit_B"},
		"Commit_B": {"Commit_A"},
		"Commit_A": {"Initial commit"},
		"Initial commit": {},
	})
	lg, err := gh.LogGraph()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, lg) {
		t.Fatalf("unexpected commit graph: %s\n", lg)
	}
}

