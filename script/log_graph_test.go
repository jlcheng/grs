package script

import (
	"reflect"
	"testing"
)

// Verifies that Scirpt.LogGraph produces a graph object
func TestLogGraph(t *testing.T) {
	const TEST_ID = "TestLogGraph"
	tmpdir, cleanup := MkTmpDir(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)

	// Creates commit log with this form
	//  a--b---c---f
	//      \     /
	//       d---e
	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "Tag_A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("tag", "Tag_B")
	gh.TouchAndCommit("C.txt", "C")
	gh.RunGit("checkout", "-b", "source_2", "Tag_B")
	gh.TouchAndCommit("D.txt", "D")
	gh.TouchAndCommit("E.txt", "E")
	gh.RunGit("checkout", "master")
	gh.RunGit("merge", "source_2", "-m", "F")

	expected := LogGraph(map[string][]string{
		"F":    {"C", "E"},
		"E":    {"D"},
		"D":    {"B"},
		"C":    {"B"},
		"B":    {"A"},
		"A":    {"init"},
		"init": {},
	})
	lg, err := gh.LogGraph()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, lg) {
		t.Fatalf("unexpected commit graph: %s\n", lg)
	}
}
