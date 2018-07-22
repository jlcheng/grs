package main

import (
	"flag"
	"jcheng/grs/grsdb"
	"fmt"
	"os"
)

func main() {
	args := grsdb.DefaultViewerOptions
	flag.BoolVar(&args.TextMode, "text", false, "treat items values as strings")
	flag.StringVar(&args.Add, "add", "", "add the specified key:value pair to the db")
	flag.Parse()

	if args.Dir = flag.Arg(0); args.Dir == "" {
		fmt.Println("dir is required\n")
		flag.Usage()
		os.Exit(1)
	}

	vopts := grsdb.DefaultViewerOptions
	vopts.Dir = args.Dir
	err := grsdb.BadgerDbAsString(args)
	if err != nil {
		fmt.Println(err)
	}
}
