package script

import (
	"jcheng/grs/shexec"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos []Repo             // a set of repositories to check and report on
	ctx   *shexec.AppContext // the application context, e.g., dependencies
	gui   AnsiGUI            // notifies the display subsystem to re-render the UI
}

func NewSyncController(repos []Repo, ctx *shexec.AppContext, gui AnsiGUI) SyncController {
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
		case BRANCH_BEHIND:
			s.AutoFFMerge()
		case BRANCH_DIVERGED:
			s.AutoRebase()
		case BRANCH_AHEAD:
			s.AutoPush()
		}

	}
}

func (d *SyncController) Run() {
	d.runIteration()
	d.gui.Run(d.repos)
}
