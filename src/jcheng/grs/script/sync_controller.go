package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos []status.Repo   // a set of repositories to check and report on
	ctx   *grs.AppContext // the application context, e.g., dependencies
	gui   AnsiGUI         // notifies the display subsystem to re-render the UI
}

func NewSyncController(repos []status.Repo, ctx *grs.AppContext, gui AnsiGUI) SyncController {
	return SyncController{
		repos: repos,
		ctx:   ctx,
		gui:   gui,
	}
}

func (d *SyncController) runIteration() {
	for i, _ := range d.repos {
		repo := &d.repos[i]
		s := NewScript(d.ctx, repo)
		s.BeforeScript()
		s.Fetch()
		s.GetRepoStatus()
		s.GetIndexStatus()

		switch repo.Branch {
		case status.BRANCH_BEHIND:
			s.AutoFFMerge()
		case status.BRANCH_DIVERGED:
			s.AutoRebase()
		}

	}
}

func (d *SyncController) Run() {
	d.runIteration()
	d.gui.Run(d.repos)
}
