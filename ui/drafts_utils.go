package ui

import (
	"jcheng/grs/script"
	"time"
)

func UpdateUI(cliUI CliUI, delay time.Duration) {
	time.Sleep(delay)
	stats := script.NewGrsStats(
		script.WithBranchstat(script.BRANCH_UPTODATE),
		script.WithIndexstat(script.INDEX_MODIFIED),
	)
	stats.CommitTime = "X seconds ago"

	cliUI.DrawGrs([]script.GrsRepo{
		script.NewGrsRepo(
			script.WithStats(stats),
			script.WithLocalGrsRepo("/foo/bar"),
			script.WithPushAllowed(true),
		),
	})
}
