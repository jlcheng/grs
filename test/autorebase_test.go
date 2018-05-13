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

func TestAutoRebase_Ok(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("^/path/to/git rev-parse", Ok(""))
	runner.AddMap("^/path/to/git rev-list", Ok("0\t0\n"))

	ctx := grs.NewAppContext()
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	repo := grs.Repo{"foo"}

	script.AutoRebase(ctx, repo, runner, rstat, false)

}

func TestAutoRebase_Test1(t *testing.T) {
	tctx := gittest.NewTestContext()

	oldwd, tmpdir := MkTmpDir(t, "AutoRebaseTest1", "TestAutoRebase_Test1")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestAutoRebase_Test1")

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
