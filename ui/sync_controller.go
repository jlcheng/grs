package ui

import (
	"jcheng/grs/script"
	"sync"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos []script.Repo      // a set of repositories to check and report on
	ctx   *script.AppContext // the application context, e.g., dependencies
	gui   AnsiGUI            // notifies the display subsystem to re-render the UI
	Cui   *CuiGUI
}

func NewSyncController(repos []script.Repo, ctx *script.AppContext, gui AnsiGUI) SyncController {
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
		script := script.NewScript(d.ctx, &d.repos[i])
		go func() {
			processRepo(script)
			wg.Done()
		}()
	}
	wg.Wait()
}

func processRepo(s *script.Script) {
	s.BeforeScript()
	s.Fetch()
	s.GetRepoStatus()
	s.GetIndexStatus()

	switch s.GetRepo().Branch {
	case script.BRANCH_BEHIND:
		s.AutoFFMerge()
	case script.BRANCH_DIVERGED:
		s.AutoRebase()
	}
	s.AutoPush()
	s.GetCommitTime()

}

func (d *SyncController) Run() {
	d.runIteration()
	// TODO hack, should properly model gui types using interface
	if d.Cui != nil {
		d.Cui.Run(d.repos)
	} else {
		d.gui.Run(d.repos)
	}
}
