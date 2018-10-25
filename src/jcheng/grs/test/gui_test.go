package test

import (
	"jcheng/grs/core"
	"jcheng/grs/script"
	"jcheng/grs/status"
	"testing"
)

// Verifies Gui starts and stops as expected
func TestGui(t *testing.T) {
	runCh := make(chan bool)
	reporter := func() []status.Repo {
		return make([]status.Repo, 0)
	}
	gui := script.NewGUI(grs.NewAppContext(), runCh, reporter, true)
	go gui.Start()
	runCh <- true
}
