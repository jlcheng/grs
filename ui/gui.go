package ui

import (
	"jcheng/grs/script"
)

type CliUI interface {
	// EventSender returns a channel one can use to poll for UI events
	EventSender() <-chan UiEvent

	// MainLoop starts the UI object
	MainLoop() error

	// DrawGrs accepts a list of repositories and renders them to screen
	DrawGrs(repo []script.GrsRepo)

	// Close will stop the UI
	Close()
}
