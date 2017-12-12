package main

import (
	"fmt"
	"jcheng/grs/gitscripts"

	"os"
)

func main() {
	repo := gitscripts.Repo{os.ExpandEnv("$HOME/github/test") }
	c, err := gitscripts.Pwd(repo)
	if err != nil {
		fmt.Println("err ", err)
		return
	}
	fmt.Println("done", c.String())
}


