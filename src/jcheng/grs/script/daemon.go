package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
	"sync"
)

// IterateRepos provides a struct that can check and report on status of a collection of repositories
type SyncDaemon struct {
	repos   []status.Repo   // a set of repositories to check and report on
	ctx     *grs.AppContext // the application context, e.g., dependencies
	done    chan struct{}   // signals SyncDaemon to shutdown
	run     chan bool       // signals SyncDaemon to run an iteration
	stopped chan struct{}   // allow other goroutines to wait for SyncDaemon to gracefully shutdown
	display chan<- bool     // notifies the display subsystem to re-render the UI
}

func NewSyncDaemon(repos []status.Repo, ctx *grs.AppContext, display chan<- bool) SyncDaemon {
	return SyncDaemon{
		repos:   repos,
		ctx:     ctx,
		done:    make(chan struct{}),
		run:     make(chan bool),
		stopped: make(chan struct{}),
		display: display,
	}
}

// StartDaemon should be invoked only once, when the SyncDaemon is started, as a goroutine
func (d *SyncDaemon) StartDaemon() {
	running := true
	var wg sync.WaitGroup // used to wait until the goroutine has been schedule
	wg.Add(1)
	go func() {
		defer close(d.stopped)
		wg.Done()
		for running {
			select {
			case <-d.run:
				grs.Debug("run signal received")
				d.runIteration()
				// non-blocking scheduling of a display event
				select {
				case d.display <- true:
				default:
				}
			case <-d.done:
				running = false
			}
		}
	}()
	wg.Wait()
}

func (d *SyncDaemon) runIteration() {
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

func (d *SyncDaemon) Run() {
	select {
	case d.run <- true:
	default:
		grs.Debug("run already scheduled")
	}
}

func (d *SyncDaemon) Shutdown() {
	close(d.done)
}

// WaitForShutdown will (cause the caller goroutine) to wait for a signal that SyncDaemon is fully stopped
func (d *SyncDaemon) WaitForShutdown() {
	<-d.stopped
}
