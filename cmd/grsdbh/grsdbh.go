package main

import (
	"flag"
	"jcheng/grs/grsdb"
	"fmt"
	"os"
)

type Args struct {
	dir       string
	asString  bool
}

func main() {

	args := Args{}
	flag.StringVar(&args.dir, "dir", "", "badgerdb to view")
	flag.BoolVar(&args.asString, "asString", false, "treat values as strings")
	flag.Parse()

	if args.dir == "" {
		fmt.Println("dir is required\n")
		os.Exit(1)
	}

	if args.asString {
		err := grsdb.BadgerDbAsString(args.dir)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("db format not specified")
		os.Exit(1)
	}

}
