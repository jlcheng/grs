package main

import (
	"jcheng/grs/ui"
	"time"
)

func main() {
	cliUI := ui.NewTviewUI()
	defer cliUI.Close()
	go ui.UpdateUI(cliUI, time.Duration(10)*time.Millisecond)
	err := cliUI.MainLoop()
	if err != nil {
		panic(err)
	}
}
