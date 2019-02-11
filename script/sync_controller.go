package script

import (
	"jcheng/grs/shexec"
	"sync"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos      []Repo             // a set of repositories to check and report on
	ctx        *shexec.AppContext // the application context, e.g., dependencies
	gui        AnsiGUI            // notifies the display subsystem to re-render the UI
	concurrent bool               // process repositories concurrently
}

func NewSyncController(repos []Repo, ctx *shexec.AppContext, gui AnsiGUI, concurrent bool) SyncController {
	return SyncController{
		repos:      repos,
		ctx:        ctx,
		gui:        gui,
		concurrent: concurrent,
	}
}

func (d *SyncController) runIteration() {
	if !d.concurrent {
		for i, _ := range d.repos {
			runSingleRepo(d.ctx, &d.repos[i])
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(len(d.repos))
		for i, _ := range d.repos {
			repo := &d.repos[i]
			go func(repo *Repo) {
				runSingleRepo(d.ctx, repo)
				wg.Done()
			}(repo)
		}
		wg.Wait()
	}
}

func runSingleRepo(ctx *shexec.AppContext, repo *Repo) {
	s := NewScript(ctx, repo)
	s.BeforeScript()
	s.Fetch()
	s.GetRepoStatus()
	s.GetIndexStatus()

	switch repo.Branch {
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
