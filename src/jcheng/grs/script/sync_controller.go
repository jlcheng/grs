package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos   []status.Repo   // a set of repositories to check and report on
	ctx     *grs.AppContext // the application context, e.g., dependencies
	display chan<- bool     // notifies the display subsystem to re-render the UI
}

func NewSyncController(repos []status.Repo, ctx *grs.AppContext, display chan<- bool) SyncController {
	return SyncController{
		repos:   repos,
		ctx:     ctx,
		display: display,
	}
}

func (d *SyncController) runIteration() {
	for _, repo := range d.repos {
		s := NewScript(d.ctx, &repo)
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
	// non-blocking scheduling of a display event
	select {
	case d.display <- true:
	default:
	}
}
