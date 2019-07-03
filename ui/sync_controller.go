package ui

import (
	"jcheng/grs/script"
	"sync"
	"time"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	grsRepos []script.GrsRepo // slice of repositories to check and report on
	Duration time.Duration    // how often to sync repos
	ui       CliUI
}

func NewSyncController(grsRepos []script.GrsRepo, ui CliUI) SyncController {
	return SyncController{
		grsRepos: grsRepos,
		ui:       ui,
	}
}

func processGrsRepo(gr *script.GrsRepo) {
	gr.ClearError()
	gr.Update()
	gr.Fetch()
	gr.UpdateRepoStatus()
	gr.UpdateIndexStatus()
	switch gr.GetStats().Branch {
	case script.BRANCH_BEHIND:
		gr.AutoFFMerge()
	case script.BRANCH_DIVERGED:
		gr.AutoRebase()
	}
	gr.AutoPush()
	gr.Update()
}

// === START: CliUI implementation ===
func (sc *SyncController) CliUIImpl() {
	done := sc.ui.DoneSender()
	ticker := time.NewTicker(sc.Duration)
	defer ticker.Stop()
	go sc.loop(done, ticker.C)
	_ = sc.ui.MainLoop()
}

// loop listens and reacts to events from the UI (SyncController.ui), including
// a clock that issues a "refresh" event every few seconds.
func (sc *SyncController) loop(done <-chan struct{}, ticker <-chan time.Time) {
	processRepoSlice := func() []script.GrsRepo {
		var wg sync.WaitGroup
		wg.Add(len(sc.grsRepos))
		for idx := range sc.grsRepos {
			go func(idx int) {
				processGrsRepo(&sc.grsRepos[idx])
				wg.Done()
			}(idx)
		}
		wg.Wait()
		return sc.grsRepos
	}
	// run at least once
	sc.ui.DrawGrs(processRepoSlice())
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
			sc.ui.DrawGrs(processRepoSlice())
		case event := <-sc.ui.EventSender():
			if event == EVENT_REFRESH {
				sc.ui.DrawGrs(processRepoSlice())
			}
		case <-done:
			break SYNC_LOOP
		}
	}
}

// === END: CliUI implementation ===
