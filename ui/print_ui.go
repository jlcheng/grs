package ui

import (
	"fmt"
	"jcheng/grs"
	"time"
)

// PrintUI is the simpler and less useful implementation of CliUI
type PrintUI struct {
	done    chan struct{}
	eventCh <-chan UiEvent
}

// NewPrintUI returns a PrintUI instance
func NewPrintUI() (*PrintUI, error) {
	return &PrintUI{
		done:    make(chan struct{}),
		eventCh: make(chan UiEvent),
	}, nil
}

func (printUI *PrintUI) MainLoop() error {
	<-printUI.done
	return nil
}

func (printUI *PrintUI) DrawGrs(repos []grs.GrsRepo) {
	fmt.Print("\033[2J\033[H")
	fmt.Println(time.Now().Format("=== Jan _2 3:04PM MST ==="))

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v, %v.\n",
			repo.GetLocal(), repo.GetStats().Branch, repo.GetStats().Index, repo.GetStats().CommitTime)
	}
}

func (printUI *PrintUI) Close() {
	close(printUI.done)
}

// EventSender returns a channel that always blocks, as the PrintUI object is too simple to generate UI events
func (printUI *PrintUI) EventSender() <-chan UiEvent {
	return printUI.eventCh
}
