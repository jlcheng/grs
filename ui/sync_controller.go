package ui

import (
	"jcheng/grs/script"
	"sync"
	"time"
)

// SyncController provides a struct that can check and report on status of a collection of repositories
type SyncController struct {
	repos []script.Repo      // a set of repositories to check and report on
	ctx   *script.AppContext // the application context, e.g., dependencies
	gui   AnsiGUI            // notifies the display subsystem to re-render the UI
	Cui   *CuiGUI
	Duration time.Duration   // how often to sync repos
	CliUI CliUI
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

// === START: Dual Loop Process ===
// RunLoops starts two goroutines 1) to sync repos and 2) tells the UI it needs to redraw. It blocks until the UI goroutine sends a quit signal.
func (d *SyncController) RunLoops() {
	quit := d.Cui.GetQuitChannel()
	ticker := time.NewTicker(d.Duration)
	defer ticker.Stop()
	syncerToUI := make(chan []script.Repo)
	defer close(syncerToUI)
	go d.uiLoop(quit, syncerToUI)
	go d.syncerLoop(quit, ticker.C, syncerToUI)
	<-quit
}

func (d *SyncController) uiLoop(quit <-chan struct{}, from <-chan []script.Repo) {
	gui := d.Cui
UI_LOOP:
	for {
		select {
		case repos := <-from:
			gui.Run(repos)
		case <-quit:
			break UI_LOOP
		}
	}
}

// repoLoop reacts to `tick` events by processes repos and sends a notification of the results to `to`
func (d *SyncController) syncerLoop(quit <-chan struct{}, ticker <-chan time.Time, syncerToUI chan<- []script.Repo) {

	processRepoSlice := func() []script.Repo {
		var wg sync.WaitGroup
		wg.Add(len(d.repos))
		for i, _ := range d.repos {
			repo := &d.repos[i]
			go func() {
				processRepo(script.NewScript(d.ctx, repo))
				wg.Done()
			}()
		}
		wg.Wait()
		return d.repos
	}

	// run at least once
	syncerToUI <- processRepoSlice()
SYNC_LOOP:
	for {
		// tie breaker in case ticker has an event and the goroutine is notified to stop q
		select {
		case <-quit:
			break SYNC_LOOP
		default:
		}

		select {
		case <-ticker:
			syncerToUI <- processRepoSlice()
		case <-quit:
			break SYNC_LOOP
		}
	}
}
// === END: Dual Loop Process ===

// === START: CliUI implementation ===
func (sc *SyncController) CliUIImpl() {
	done := sc.CliUI.DoneSender()
	ticker := time.NewTicker(sc.Duration)
	defer ticker.Stop()
	syncerToUI := make(chan []script.Repo)
	defer close(syncerToUI)
	go sc.uiDispatchLoop(done, syncerToUI)
	go sc.appLoop(done, ticker.C, syncerToUI)
	sc.CliUI.MainLoop()
}

func (sc *SyncController) uiDispatchLoop(done <-chan struct{}, from <-chan []script.Repo) {
UI_LOOP:
	for {
		select {
		case repos := <-from:
			sc.CliUI.Draw(repos)
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