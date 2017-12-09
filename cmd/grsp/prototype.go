package main

import (
	"fmt"
	"jcheng/grs/gitscripts"

)

func main() {
	c, err := gitscripts.Status()
	if err != nil {
		fmt.Println("err ", err)
		return
	}
	fmt.Println("done", c)

}


