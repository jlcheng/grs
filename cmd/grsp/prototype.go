package main

import (
	"fmt"
	"jcheng/grs/gitscripts"

)

func main() {
	c, err := gitscripts.Status()
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(c.Stdout)

}


