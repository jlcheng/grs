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
	syncerToUI := make(chan []script.GrsRepo)
	defer close(syncerToUI)
	//go sc.uiDispatchLoop(done, syncerToUI)
	//go sc.appLoop(done, ticker.C, syncerToUI)
	go sc.uiGrsDispatchLoop(done, syncerToUI)
	go sc.appGrsLoop(done, ticker.C, syncerToUI)
	_ = sc.ui.MainLoop()
}

func (sc *SyncController) uiGrsDispatchLoop(done <-chan struct{}, from <-chan []script.GrsRepo) {
UI_LOOP:
	for {
		select {
		case repos := <-from:
			sc.ui.DrawGrs(repos)
		case <-done:
			break UI_LOOP
		}
	}
}

func (sc *SyncController) appGrsLoop(done <-chan struct{}, ticker <-chan time.Time, syncerToUI chan<- []script.GrsRepo) {
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
		case event := <-sc.ui.EventSender():
			if event == EVENT_REFRESH {
				syncerToUI <- processRepoSlice()
			}
		case <-done:
			break SYNC_LOOP
		}
	}
}

// === END: CliUI implementation ===
