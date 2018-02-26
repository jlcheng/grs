package test

import (
	"testing"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"jcheng/grs/grs"
	"jcheng/grs/grsdb"
	"time"
)

func TestFetch_Git_Fail(t *testing.T) {
	runner := NewMockRunner()
	runner.Add(Error("failed"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	ctx := grs.NewAppContext()
	script.Fetch(ctx, runner, rstat, grs.Repo{Path:"/repo"})
	if rstat.Branch != status.BRANCH_UNKNOWN {
		t.Errorf("expected %s, got: %v\n", status.BRANCH_UNKNOWN, rstat.Branch)
	}
}

func TestFetch_Git_OK(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	ctx := grs.NewAppContext()
	script.Fetch(ctx, runner, rstat, grs.Repo{Path:"/repo"})
	if rstat.Dir == status.DIR_INVALID {
		t.Error("Unexpected rstat.Dir, got DIR_INVALID")
	}
}

func TestFetch_Modified_Update(t *testing.T) {
	runner := NewMockRunner()
	runner.AddMap("git", Ok("0"))
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	ctx := grs.NewAppContext()
	script.Fetch(ctx, runner, rstat, grs.Repo{Path:"/repo"})
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
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	ctx := grs.NewAppContext()
	ctx.DB().Repos = append(ctx.DB().Repos, grsdb.Repo{Id:"/repo", FetchedSec:1})
	script.Fetch(ctx, runner, rstat, grs.Repo{Path:"/repo"})
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
	rstat := status.NewRStat()
	rstat.Dir = status.DIR_VALID
	ctx := grs.NewAppContext()
	fetchTime := time.Now().Unix()
	ctx.DB().Repos = append(ctx.DB().Repos, grsdb.Repo{Id:"/repo", FetchedSec:fetchTime})
	script.Fetch(ctx, runner, rstat, grs.Repo{Path:"/repo"})
	db := ctx.DB()
	if l := len(db.Repos); l != 1 {
		t.Errorf("Expected len(db.Repos) == 1, got %v\n", l)
	}
	dbrepo := ctx.DB().Repos[0]
	if dbrepo.FetchedSec != fetchTime {
		t.Errorf("Expected dbrepo.FetchedSec != 1, got %v\n", dbrepo.FetchedSec)
	}
}

