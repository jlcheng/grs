package script

import (
	"fmt"
	"jcheng/grs/core"
	"jcheng/grs/status"
	"sync"
)

type Reporter func() []status.Repo

type AnsiGUI struct {
	ctx      *grs.AppContext // the application context, e.g., dependencies
	done     chan struct{}   // singals GUI to shutdown
	run      <-chan bool     // signals GUI to render state
	stopped  chan struct{}   // allow other goroutines to wait for GUI to gracefully shutdown
	reporter Reporter        // provides status on repos
}

func NewGUI(ctx *grs.AppContext, run <-chan bool, reporter Reporter) AnsiGUI {
	return AnsiGUI{
		ctx:      ctx,
		done:     make(chan struct{}),
		run:      run,
		stopped:  make(chan struct{}),
		reporter: reporter,
	}
}

func (gui *AnsiGUI) Start() {
	running := true
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(gui.stopped)
		wg.Done()
		for running {
			select {
			case <-gui.run:
				gui.runIteration()
			case <-gui.done:
				running = false
			}
		}
	}()
	wg.Wait()
}

func (gui *AnsiGUI) runIteration() {
	for _, repo := range gui.reporter() {
		fmt.Printf("repo [%v] status IS %v, %v.\n",
				repo.Path, colorB(repo.Branch), colorI(repo.Index))
	}
}

func (gui *AnsiGUI) Shutdown() {
	close(gui.done)
}

func (gui *AnsiGUI) WaitShutdown() {
	<-gui.stopped
}


func colorI(s status.Indexstat) string {
	if s == status.INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorB(s status.Branchstat) string {
	if s == status.BRANCH_UPTODATE {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}