package ui

//go:generate stringer -type UiEvent -output gui_event_strings.go

// UiEvent models a enumeration UI events supported by the SyncController
type UiEvent int

const (
	EVENT_REFRESH UiEvent = iota
	EVENT_QUIT
)
