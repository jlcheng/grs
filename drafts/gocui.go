package main

import (
	"jcheng/grs/ui"
	"jcheng/grs/script"	
	"log"
	"time"
)

func main() {
	ui, err := ui.NewConsoleUI()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		ui.DrawGrs([]script.GrsRepo {
			script.NewGrsRepo(
				script.WithLocalGrsRepo("/foo/bar"),
				script.WithPushAllowed(true),
			),
		})
	}()
	
	err = ui.MainLoop()
	if err != nil {
		log.Println(err)
	}
}
