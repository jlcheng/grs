package script

import (
	"fmt"
	"jcheng/grs/core"
	"jcheng/grs/status"
)

type Reporter func() []status.Repo

type AnsiGUI struct {
	ctx      *grs.AppContext // the application context, e.g., dependencies
	run      <-chan bool     // signals GUI to render state
	stopped  chan struct{}   // notifies outside world that we're done
	reporter Reporter        // provides status on repos
	clr      bool            // if true, clears screen before each iteration
}

func NewGUI(ctx *grs.AppContext, run <-chan bool, reporter Reporter, clr bool) AnsiGUI {
	return AnsiGUI{
		ctx:      ctx,
		run:      run,
		stopped:  make(chan struct{}),
		reporter: reporter,
		clr:      clr,
	}
}

// Start will
func (gui *AnsiGUI) Start() {
	for {
		_, open := <-gui.run
		if !open {
			close(gui.stopped)
			return
		}
		gui.runIteration()
	}
}

func (gui *AnsiGUI) runIteration() {
	// setup/clear screen
	if gui.clr {
		fmt.Print("\033[2J\033[H")
	}

	for _, repo := range gui.reporter() {
		fmt.Printf("repo [%v] status IS %v, %v.\n",
				repo.Path, colorB(repo.Branch), colorI(repo.Index))
	}
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