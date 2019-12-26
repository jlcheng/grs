package ui

import (
	"jcheng/grs"
	"jcheng/grs/base"
	"sync"
	"time"
)

// SyncController reports on the status of repositories
type SyncController struct {
	// Slice of repositories to check and report on
	grsRepos []grs.GrsRepo

	// How often to sync repos
	duration time.Duration

	// Used to render the status of repositories
	ui CliUI
}

// ControllerEvent describes an event dispatched within SyncController's framework
type ControllerEvent struct {
	Type  UiEvent
	Repos []grs.GrsRepo
}

// NewSyncController allocates a SyncController struct with the given list of repos and UI
func NewSyncController(grsRepos []grs.GrsRepo, ui CliUI, duration time.Duration) SyncController {
	return SyncController{
		grsRepos: grsRepos,
		ui:       ui,
		duration: duration,
	}
}

// processGrsRepo describes the routine for synchronizing a repository
func processGrsRepo(gr *grs.GrsRepo) {
	gr.ClearError()
	gr.Update()
	gr.Fetch()
	gr.UpdateRepoStatus()
	gr.UpdateIndexStatus()
	switch gr.GetStats().Branch {
	case grs.BRANCH_BEHIND:
		gr.AutoFFMerge()
	case grs.BRANCH_DIVERGED:
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
	case EVENT_QUIT:
		sc.ui.Close()
	default:
		base.DebugFull("", "SyncCtrl", "unexpected event: %v", event.Type)
	}
}

// refresh calls `processGrsRepo` against each repo concurrently
func (sc *SyncController) refresh() []grs.GrsRepo {
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
// The event dispatch loop will handle events from the UI, such as
// "refresh" and "quit".
//
// The SyncController will also start a ticker that emits"refresh"
// on a regular interval. This interval is configured by the
// controller's 'duration' property.
func (sc *SyncController) loop() {
	ticker := time.NewTicker(sc.duration)
	defer ticker.Stop()

	// run at least once
	sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
	for {
		select {
		case <-ticker.C:
			sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
		case event := <-sc.ui.EventSender():
			if event == EVENT_REFRESH {
				sc.OnEvent(ControllerEvent{Type: EVENT_REFRESH, Repos: sc.refresh()})
			}
			if event == EVENT_QUIT {
				sc.OnEvent(ControllerEvent{Type: EVENT_QUIT, Repos: nil})
			}
		}
	}
}
