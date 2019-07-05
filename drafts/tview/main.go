package main

import (
	"jcheng/grs/ui"
)

func main() {
	cliUI := ui.NewTviewUI()
	err := cliUI.MainLoop()
	if err != nil {
		panic(err)
	}
}
