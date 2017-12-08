package main

import (
	"fmt"
	"gitscripts"

)

func main() {
	c := gitscripts.Status()
	fmt.Print(c.Stdout)

}


