package test

import (
	"io/ioutil"
	"jcheng/grs/gittest"
	"jcheng/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"testing"
)

/*
cloned_repo/master rebases without conflicts on to @{UPSTREAM}

a--b---c---f  @{UPSTREAM} origin/master
 \  \     /
  \  d---e    origin/branch_B
   \
    g---h     cloned_repo/master
*/
func TestAutoRebase_IT_Test_2(t *testing.T) {
	exec := gittest.NewExecRunner()

	oldwd, tmpdir := MkTmpDir(t, "AutoRebase_IT_Test_2", "TestAutoRebase_IT_Test_2")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_IT_Test_2")
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

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
	exec.TouchAndCommit("d.txt", "Commit_D")
	exec.TouchAndCommit("e.txt", "Commit_E")
	exec.Exec(git, "checkout", "master")
	exec.Exec(git, "merge", "--no-ff", "branch_B")

	exec.Chdir("..")
	exec.Chdir("dest")
	exec.TouchAndCommit("g.txt", "Commit_G")
	exec.TouchAndCommit("h.txt", "Commit_H")

	if exec.Err() != nil {
		t.Fatal("test setup failed")
	}

	ctx := grs.NewAppContextWithRunner(exec.Runner())
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	s := script.NewScript(ctx, repo)
	s.Fetch()
	s.AutoRebase()
	s.GetRepoStatus()

	if repo.Branch != status.BRANCH_AHEAD {
		t.Fatalf("expected BRANCH_UPTODATE, but was %v\n", repo.Branch)
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
	exec := gittest.NewExecRunner()

	oldwd, tmpdir := MkTmpDir(t, "AutoRebase_IT_Test_3", "TestAutoRebase_IT_Test_3")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_IT_Test_3")
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

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

	ctx := grs.NewAppContextWithRunner(exec.Runner())
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	s := script.NewScript(ctx, repo)
	s.Fetch()
	s.AutoRebase()
	s.GetRepoStatus()
	if repo.Branch != status.BRANCH_DIVERGED {
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

func CleanTmpDir(t *testing.T, oldwd string, tmpdir string, errid string) {

	if err := os.Chdir(oldwd); err != nil {
		t.Fatal(errid, err)
	}
	if err := os.RemoveAll(tmpdir); err != nil {
		t.Fatal(errid, err)
	}
}
