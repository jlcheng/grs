package gitscripts

import (
	"os/exec"
	"bytes"
)

func Status() *exec.Cmd {
	c := exec.Command("git status")
	out := new(bytes.Buffer)
	c.Stdout = out
	c.Run()
	return c
}
