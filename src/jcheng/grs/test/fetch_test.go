package test

import (
	"jcheng/grs/core"
	"jcheng/grs/grsdb"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
	"time"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, repo.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	if repo.Dir == status.DIR_INVALID {
		t.Error("Unexpected repo.Dir, got DIR_INVALID")
	}
}

func TestFetch_Modified_Update(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := status.NewRepo("")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	script.NewScript(ctx, repo).Fetch()
	db := ctx.DB()
	if l := len(db.Repos); l != 1 {
		t.Errorf("Expected len(db.Repos) == 1, got %v\n", l)
	}
	dbrepo := ctx.DB().Repos[0]
	if dbrepo.FetchedSec == 0 {
		t.Error("Expected dbrepo.FetchedSec != 0, got 0")
	}
}

func TestFetch_Modified_Update_Existing(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := status.NewRepo("/repo")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	ctx.DB().Repos = append(ctx.DB().Repos, grsdb.RepoDTO{Id: repo.Path, FetchedSec: 1})
	script.NewScript(ctx, repo).Fetch()
	db := ctx.DB()
	if l := len(db.Repos); l != 1 {
		t.Errorf("Expected len(db.Repos) == 1, got %v\n", l)
	}
	dbrepo := ctx.DB().Repos[0]
	if dbrepo.FetchedSec == 1 {
		t.Errorf("Expected dbrepo.FetchedSec != 1, got %v\n", dbrepo.FetchedSec)
	}
}

func TestFetch_Modified_Update_NOP(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	repo := status.NewRepo("foo")
	repo.Dir = status.DIR_VALID
	ctx := grs.NewAppContextWithRunner(runner)
	fetchTime := time.Now().Unix()
	ctx.DB().Repos = append(ctx.DB().Repos, grsdb.RepoDTO{Id: repo.Path, FetchedSec: fetchTime})
	script.NewScript(ctx, repo).Fetch()
	db := ctx.DB()

	if l := len(db.Repos); l != 1 {
		t.Errorf("Expected len(db.Repos) == 1, got %v\n", l)
	}
	dbrepo := ctx.DB().Repos[0]
	if dbrepo.FetchedSec != fetchTime {
		t.Errorf("Expected dbrepo.FetchedSec != 1, got %v\n", dbrepo.FetchedSec)
	}
}
