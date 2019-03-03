package script

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// Verifies we can rebase without conflicts on top of source
/*
Given

    a--b---c---f
     \  \     /   remote branch, which is a non-trivial graph
      \  d---e
       \
        g---h     local branch, which is a trivial graph of a->h->i

Then AutoRebase() should create

    a--b---c---f
        \     / \
         d---e   g---h   local branch contains all changes from source
*/
func TestAutoRebase_IT_Test_2(t *testing.T) {
	const TEST_ID = "TestAutoRebase_IT_Test_2"
	tmpdir, cleanup := MkTmpDir1(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	repo := NewRepo(gh.Getwd())
	repo.PushAllowed = true
	s := NewScript(NewAppContext(), repo)
	s.BeforeScript()

	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("tag", "B")
	gh.TouchAndCommit("C.txt", "C")
	gh.RunGit("checkout", "-b", "source_2", "B")
	gh.TouchAndCommit("D.txt", "D")
	gh.TouchAndCommit("E.txt", "E")
	gh.RunGit("checkout", "master")
	gh.RunGit("merge", "source_2", "-m", "F")
	gh.RunGit("push")
	gh.RunGit("reset", "--hard", "A")
	gh.TouchAndCommit("G.txt", "G")
	gh.TouchAndCommit("H.txt", "H")

	s.AutoRebase()
	s.Update()

	got := LogGraph(map[string][]string{
		"A":    {"init"},
		"B":    {"A"},
		"C":    {"B"},
		"D":    {"B"},
		"E":    {"D"},
		"F":    {"C", "E"},
		"G":    {"F"},
		"H":    {"G"},
		"init": {},
	})
	expected, _ := gh.LogGraph()
	if !reflect.DeepEqual(got, expected) {
		t.Fatal("unexpected commit graph", got)
	}
}

// Verifies that rebase does not happen when there is a conflict upstream
/*
Given the following, where h and c conflicts

a--b---c---f
 \  \     /   remote branch, where c modifies conflicts.txt
  \  d---e
   \
    g---h     local branch, where h has conflicting changes with c

Then AutoRebase() should create in the local branch

    a--b
        \
         d---e---g---h  local branch gets d and e from source, but does not have c and f, which contains conflicts
*/
func TestAutoRebase_IT_Test_3(t *testing.T) {
	const TEST_ID = "TestAutoRebase_IT_Test_3"
	tmpdir, cleanup := MkTmpDir1(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	repo := NewRepo(gh.Getwd())
	repo.PushAllowed = true
	s := NewScript(NewAppContext(), repo)
	s.BeforeScript()

	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("tag", "B")
	gh.SetContents("conflict.txt", "C")
	gh.Add("conflict.txt")
	gh.TouchAndCommit("C.txt", "C")
	gh.RunGit("checkout", "-b", "source_2", "B")
	gh.TouchAndCommit("D.txt", "D")
	gh.TouchAndCommit("E.txt", "E")
	gh.RunGit("checkout", "master")
	gh.RunGit("merge", "source_2", "-m", "F")
	gh.RunGit("push")
	gh.RunGit("reset", "--hard", "A")

	gh.TouchAndCommit("G.txt", "G")
	gh.SetContents("conflict.txt", "H")
	gh.Add("conflict.txt")
	gh.TouchAndCommit("H.txt", "H")

	if gh.Err() != nil {
		t.Fatal("test setup failed", gh.Err())
	}

	s.Fetch()
	s.AutoRebase()
	s.GetRepoStatus()

	got := LogGraph(map[string][]string{
		"A":    {"init"},
		"B":    {"A"},
		"D":    {"B"},
		"E":    {"D"},
		"G":    {"E"},
		"H":    {"G"},
		"init": {},
	})
	expected, _ := gh.LogGraph()
	if !reflect.DeepEqual(got, expected) {
		t.Fatal("unexpected commit graph", got)
	}
}

func MkTmpDir(t *testing.T, prefix string, errid string) (oldwd string, d string) {
	var err error
	oldwd, err = os.Getwd()
	if err != nil {
		t.Fatal(errid, err)
	}
	d, err = ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(errid, err)
	}
	return oldwd, d
}

// MkTmpDir creates a temporary directory usiing ioutil.TempDir and calls t.Fatal if the attempt fails. On success, it
// returns:
// - the created directory
// - a no-arg function which deletes the temp directory and os.Chdir to the current working directory
func MkTmpDir1(t *testing.T, errid string) (string, func()) {
	var err error
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(errid, err)
	}
	tempDir, err := ioutil.TempDir("", errid)
	if err != nil {
		t.Fatal(errid, err)
	}

	return tempDir, func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal(errid, err)
		}
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(errid, err)
		}
	}
}


func CleanTmpDir(t *testing.T, oldwd string, tmpdir string, errid string) {

	if err := os.Chdir(oldwd); err != nil {
		t.Fatal(errid, err)
	}
	if err := os.RemoveAll(tmpdir); err != nil {
		t.Fatal(errid, err)
	}
}
