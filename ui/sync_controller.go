package ui

import (
	"jcheng/grs/base"
	"jcheng/grs/script"
	"sync"
	"time"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	grsRepos []script.GrsRepo // slice of repositories to check and report on
	duration time.Duration    // how often to sync repos
	ui       CliUI
}

// ControllerEvent describes an event dispatched within SyncController's framework
type ControllerEvent struct {
	Type  UiEvent
	Repos []script.GrsRepo
}

// NewSyncController allocates a SyncController struct with the given list of repos and UI
func NewSyncController(grsRepos []script.GrsRepo, ui CliUI, duration time.Duration) SyncController {
	return SyncController{
		grsRepos: grsRepos,
		ui:       ui,
		duration: duration,
	}
}

// processGrsRepo describes the routine for synchronizing a repository
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

// Run starts this controller and blocks until the UI receives a 'Done' signal.
func (sc *SyncController) Run() {
	go sc.loop()
	_ = sc.ui.MainLoop()
}

// OnEvent processes an event within SyncController's framework
func (sc *SyncController) OnEvent(event ControllerEvent) {
	switch event.Type {
	case EVENT_REFRESH:
		sc.ui.DrawGrs(event.Repos)
	default:
		base.Debug("unexpected event: %v", event.Type)
	}
}

// refresh calls `processGrsRepo` against each repo concurrently
func (sc *SyncController) refresh() []script.GrsRepo {
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

// loop starts the event dispatch loop.
//
// The event dispatch loop will handle events from the UI, for
// example, an explict 'refresh' event.  It also handles the 'Done'
// event and stops the dispatch loop.
//
// The SyncController will also start a ticker that emits a "refresh"
// on a regular interval. This interval is configured by the
// controller's 'duration' property.
func (sc *SyncController) loop() {
	ticker := time.NewTicker(sc.duration)
	defer ticker.Stop()

	// run at least once
	sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
SYNC_LOOP:
	for {
		// tie breaker in case ticker has an event and the goroutine is notified to stop q
		select {
		case <-sc.ui.DoneSender():
			break SYNC_LOOP
		default:
		}

		select {
		case <-ticker.C:
			sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
		case event := <-sc.ui.EventSender():
			if event == EVENT_REFRESH {
				sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
			}
		case <-sc.ui.DoneSender():
			break SYNC_LOOP
		}
	}
}
