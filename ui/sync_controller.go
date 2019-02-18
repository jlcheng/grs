package ui

import (
	"jcheng/grs/script"
	"sync"
	"time"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos    []script.Repo      // a set of repositories to check and report on
	ctx      *script.AppContext // the application context, e.g., dependencies
	Duration time.Duration   // how often to sync repos
	ui       CliUI
}

func NewSyncController(repos []script.Repo, ctx *script.AppContext, ui CliUI) SyncController {
	return SyncController{
		repos: repos,
		ctx:   ctx,
		ui:    ui,
	}
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


// === START: CliUI implementation ===
func (sc *SyncController) CliUIImpl() {
	done := sc.ui.DoneSender()
	ticker := time.NewTicker(sc.Duration)
	defer ticker.Stop()
	syncerToUI := make(chan []script.Repo)
	defer close(syncerToUI)
	go sc.uiDispatchLoop(done, syncerToUI)
	go sc.appLoop(done, ticker.C, syncerToUI)
	sc.ui.MainLoop()
}

func (sc *SyncController) uiDispatchLoop(done <-chan struct{}, from <-chan []script.Repo) {
UI_LOOP:
	for {
		select {
		case repos := <-from:
			sc.ui.Draw(repos)
		case <-done:
			break UI_LOOP
		}
	}
}

func (sc *SyncController) appLoop(done <-chan struct{}, ticker <-chan time.Time, syncerToUI chan<- []script.Repo) {
	processRepoSlice := func() []script.Repo {
		var wg sync.WaitGroup
		wg.Add(len(sc.repos))
		for i, _ := range sc.repos {
			repo := &sc.repos[i]
			go func() {
				processRepo(script.NewScript(sc.ctx, repo))
				wg.Done()
			}()
		}
		wg.Wait()
		return sc.repos
	}
	// run at least once
	syncerToUI <- processRepoSlice()
SYNC_LOOP:
	for {
		// tie breaker in case ticker has an event and the goroutine is notified to stop q
		select {
		case <-done:
			break SYNC_LOOP
		default:
		}

		select {
		case <-ticker:
			syncerToUI <- processRepoSlice()
		case <-done:
			break SYNC_LOOP
		}
	}
}
// === END: CliUI implementation ===