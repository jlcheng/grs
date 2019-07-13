package ui

import (
	"jcheng/grs/script"
	"time"
	"fmt"
)

func UpdateUI(cliUI CliUI, delay time.Duration) {
	time.Sleep(delay)
	repos := make([]script.GrsRepo, 0)
	stats := script.NewGrsStats(
		script.WithBranchstat(script.BRANCH_UPTODATE),
		script.WithIndexstat(script.INDEX_MODIFIED),
	)
	stats.CommitTime = "X seconds ago"
	repos = append(repos, 
		script.NewGrsRepo(
			script.WithStats(stats),
			script.WithLocalGrsRepo("/foo/bar"),
			script.WithPushAllowed(true),
		))
	
	stats = script.NewGrsStats(
		script.WithBranchstat(script.BRANCH_UPTODATE),
		script.WithIndexstat(script.INDEX_UNMODIFIED),
	)
	stats.CommitTime = "X seconds ago"	
	repos = append(repos,
		script.NewGrsRepo(
			script.WithStats(stats),
			script.WithLocalGrsRepo("/foo/repo2"),
			script.WithPushAllowed(false),
		))

	stats = script.NewGrsStats(
		script.WithBranchstat(script.BRANCH_UPTODATE),
		script.WithIndexstat(script.INDEX_UNMODIFIED),
	)
	stats.CommitTime = "Z minutes ago"
	repos = append(repos,
		script.NewGrsRepo(
			script.WithStats(stats),
			script.WithLocalGrsRepo("/foo/repo3"),
			script.WithPushAllowed(false),
			script.WithError(fmt.Errorf("foo")),
		))
	cliUI.DrawGrs(repos)
	
}
