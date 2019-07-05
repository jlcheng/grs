package main

import (
	"jcheng/grs/ui"
	"log"
	"time"
)

func main() {
	cliUI, err := ui.NewConsoleUI()
	if err != nil {
		log.Fatal(err)
	}
	defer cliUI.Close()
	go ui.UpdateUI(cliUI, time.Duration(10) * time.Millisecond)
	
	err = cliUI.MainLoop()
	if err != nil {
		log.Println(err)
	}
}
