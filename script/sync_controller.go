package script

import (
	"sync"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos []Repo             // a set of repositories to check and report on
	ctx   *AppContext // the application context, e.g., dependencies
	gui   AnsiGUI            // notifies the display subsystem to re-render the UI
}

func NewSyncController(repos []Repo, ctx *AppContext, gui AnsiGUI) SyncController {
	return SyncController{
		repos: repos,
		ctx:   ctx,
		gui:   gui,
	}
}

func (d *SyncController) runIteration() {
	var wg sync.WaitGroup
	wg.Add(len(d.repos))
	for i, _ := range d.repos {
		script := NewScript(d.ctx, &d.repos[i])
		go func() {
			processRepo(script)
			wg.Done()
		}()
	}
	wg.Wait()
}

func processRepo(s *Script) {
	s.BeforeScript()
	s.Fetch()
	s.GetRepoStatus()
	s.GetIndexStatus()

	switch s.repo.Branch {
	case BRANCH_BEHIND:
		s.AutoFFMerge()
	case BRANCH_DIVERGED:
		s.AutoRebase()
	}
	s.AutoPush()
	s.GetCommitTime()

}

func (d *SyncController) Run() {
	d.runIteration()
	d.gui.Run(d.repos)
}
