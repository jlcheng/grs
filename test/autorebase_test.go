package test

import (
	"io/ioutil"
	"jcheng/grs/config"
	"jcheng/grs/gittest"
	"jcheng/grs/grs"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"os"
	"strings"
	"testing"
	"fmt"
)

func TestToClonePath(t *testing.T) {
	a := string(os.PathSeparator) + "foo"
	b := "foo"
	aclone := script.ToClonePath(a)
	bclone := script.ToClonePath(b)
	if aclone == bclone {
		t.Fatal("a and b must not yield same 'clone path'", aclone, bclone)
	}
	if !strings.HasPrefix(aclone, config.UserPrefDir) {
		t.Fatal("The 'clone path' must start with user pref directory")
	}
}

func TestAutoRebase_IT_Test_1(t *testing.T) {
	tctx := gittest.NewTestContext()

	oldwd, tmpdir := MkTmpDir(t, "AutoRebase_IT_Test_1", "TestAutoRebase_IT_Test_1")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_IT_Test_1")

	if err := gittest.InitTest1(tctx, tmpdir); err != nil {
		t.Fatal(err, "TestAutoRebase_Test1")
	}

	ctx := grs.NewAppContext()
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	repo := grs.Repo{"foo"}
	runner := tctx.GetRunner()

	script.AutoRebase(ctx, repo, runner, rstat, false)

}

/*
cloned_repo/master rebases without conflicts on to @{UPSTREAM}

a--b---c---f  @{UPSTREAM} origin/master
 \  \     /
  \  d---e    origin/branch_B
   \
    g---h     cloned_repo/master
*/
func TestAutoRebase_IT_Test_2(t *testing.T) {
	tctx := gittest.NewTestContext()

	oldwd, tmpdir := MkTmpDir(t, "AutoRebase_IT_Test_2", "TestAutoRebase_IT_Test_2")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_IT_Test_2")

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if ok {
				t.Fatal(err)
			}
		}
	}()

	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

	git := tctx.Git()
	tctx.Mkdir("source")
	tctx.Chdir("source")
	tctx.Exec(git, "init")
	tctx.TouchAndCommit(".gitignore", "Commit_A")
	tctx.Chdir("..")
	tctx.Exec(git, "clone", "source", "dest")

	tctx.Chdir("./source")
	tctx.TouchAndCommit("b.txt", "Commit_B")
	tctx.TouchAndCommit("c.txt", "Commit_C")
	tctx.Exec(git, "checkout", "-b", "branch_B")
	tctx.TouchAndCommit("d.txt", "Commit_D")
	tctx.TouchAndCommit("e.txt", "Commit_E")
	tctx.Exec(git, "checkout", "master")
	tctx.Exec(git, "merge", "--no-ff", "branch_B")

	tctx.Chdir("..")
	tctx.Chdir("dest")
	tctx.TouchAndCommit("g.txt", "Commit_G")
	tctx.TouchAndCommit("h.txt", "Commit_H")

	ctx := grs.NewAppContext()
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	repo := grs.Repo{"foo"}
	runner := tctx.GetRunner()
	script.Fetch(ctx, runner, rstat, repo)

	err := script.AutoRebase(ctx, repo, runner, rstat, false)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
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

