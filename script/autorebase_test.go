package script

import (
	"io/ioutil"
	"os"
	"testing"
)

// dest rebases without conflicts on top of source
/*
Given

    a--b---c---e
     \  \     /   source, which is a non-trivial graph
      \  d---e
       \
        g---h     dest, which is a trivial graph of a->g->h

Then AutoRebase() should create

    a--b---c---e
        \     / \
         d---e   g---h dest, which is a non-trivial graph of a->(complex)->f->g->h
*/
func TestAutoRebase_IT_Test_2(t *testing.T) {
	const TEST_ID = "TestAutoRebase_IT_Test_2"
	tmpdir, cleanup := MkTmpDir1(t, TEST_ID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(true), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	repo := NewRepo(gh.Getwd())
	repo.PushAllowed = true
	s := NewScript(NewAppContext(), repo)
	s.BeforeScript()

	gh.TouchAndCommit("A.txt", "Commit_A")
	gh.GitExec("tag", "Commit_A")
	gh.TouchAndCommit("B.txt", "Commit_B")
	gh.GitExec("tag", "Commit_B")
	gh.TouchAndCommit("C.txt", "Commit_C")
	gh.GitExec("checkout", "-b", "source_2", "Commit_B")
	gh.TouchAndCommit("D.txt", "Commit_D")
	gh.TouchAndCommit("E.txt", "Commit_E")
	gh.GitExec("checkout", "master")
	gh.GitExec("merge", "source_2")
	gh.GitExec("push")
	gh.GitExec("reset", "--hard", "Commit_A")
	gh.TouchAndCommit("G.txt", "Commit_G")
	gh.TouchAndCommit("H.txt", "Commit_H")

	s.AutoRebase()
	s.Update()

	gh.GitExec("log", "--pretty=%h %p") // TODO kept for debugging
	gh.GitExec("log", "--pretty=%h %s") // TODO kept for debugging
	if repo.Index != INDEX_UNMODIFIED || repo.Branch != BRANCH_AHEAD {
		t.Fatal("repo is not uptodate and unmodified")
	}
}

/*
cloned_repo/master rebases with conflicts on to @{UPSTREAM}

a--b---c---f  @{UPSTREAM} origin/master
 \  \     /
  \  d---e    origin/branch_B
   \
    g---h     cloned_repo/master (g has a conflict with commit d)
*/
func TestAutoRebase_IT_Test_3(t *testing.T) {
	oldwd, tmpdir := MkTmpDir(t, "AutoRebase_IT_Test_3", "TestAutoRebase_IT_Test_3")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_IT_Test_3")
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

	exec := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	git := exec.Git()
	exec.Mkdir("source")
	exec.Chdir("source")
	exec.Exec(git, "init")
	exec.TouchAndCommit(".gitignore", "Commit_A")
	exec.Chdir("..")
	exec.Exec(git, "clone", "source", "dest")

	exec.Chdir("./source")
	exec.TouchAndCommit("b.txt", "Commit_B")
	exec.TouchAndCommit("c.txt", "Commit_C")
	exec.Exec(git, "checkout", "-b", "branch_B")
	exec.SetContents("conflict.txt", "D")
	exec.Add("conflict.txt")
	exec.TouchAndCommit("d.txt", "Commit_D")
	exec.TouchAndCommit("e.txt", "Commit_E")
	exec.Exec(git, "checkout", "master")
	exec.Exec(git, "merge", "--no-ff", "branch_B")

	exec.Chdir("..")
	exec.Chdir("dest")
	exec.SetContents("conflict.txt", "G")

	exec.Add("conflict.txt")
	exec.TouchAndCommit("g.txt", "Commit_G")
	exec.TouchAndCommit("h.txt", "Commit_H")

	if exec.Err() != nil {
		t.Fatal("test setup failed", exec.Err())
	}

	ctx := NewAppContext(WithCommandRunner(exec.CommandRunner()))
	repo := NewRepo(exec.Getwd())
	repo.Dir = DIR_VALID
	s := NewScript(ctx, repo)
	s.Fetch()
	s.AutoRebase()
	s.GetRepoStatus()
	if repo.Branch != BRANCH_DIVERGED {
		t.Fatalf("expected BRANCH_DIVERGED, but was %v\n", repo.Branch)
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
