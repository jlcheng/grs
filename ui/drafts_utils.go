package ui

import (
	"fmt"
	"jcheng/grs"
	"time"
)

func UpdateUI(cliUI CliUI, delay time.Duration) {
	time.Sleep(delay)
	repos := make([]grs.GrsRepo, 0)
	stats := grs.NewGrsStats(
		grs.WithBranchstat(grs.BRANCH_UPTODATE),
		grs.WithIndexstat(grs.INDEX_MODIFIED),
	)
	stats.CommitTime = "X seconds ago"
	repos = append(repos,
		grs.NewGrsRepo(
			grs.WithStats(stats),
			grs.WithLocalGrsRepo("/foo/bar"),
			grs.WithPushAllowed(true),
		))

	stats = grs.NewGrsStats(
		grs.WithBranchstat(grs.BRANCH_UPTODATE),
		grs.WithIndexstat(grs.INDEX_UNMODIFIED),
	)
	stats.CommitTime = "X seconds ago"
	repos = append(repos,
		grs.NewGrsRepo(
			grs.WithStats(stats),
			grs.WithLocalGrsRepo("/foo/repo2"),
			grs.WithPushAllowed(false),
		))

	stats = grs.NewGrsStats(
		grs.WithBranchstat(grs.BRANCH_UPTODATE),
		grs.WithIndexstat(grs.INDEX_UNMODIFIED),
	)
	stats.CommitTime = "Z minutes ago"
	repos = append(repos,
		grs.NewGrsRepo(
			grs.WithStats(stats),
			grs.WithLocalGrsRepo("/foo/repo3"),
			grs.WithPushAllowed(false),
			grs.WithError(fmt.Errorf("foo")),
		))
	cliUI.DrawGrs(repos)

}
