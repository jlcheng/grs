package main

import (
	"fmt"
	"jcheng/grs/grsio"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: grscpdir <src> <dest>")
		os.Exit(1)
	}
	src := os.Args[1]
	dest := os.Args[2]
	err := grsio.CopyDir(src, dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
