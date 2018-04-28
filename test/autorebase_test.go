package test

import (
	"jcheng/grs/config"
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
	script.AutoRebase(ctx, runner, rstat, repo)

}
