package gitscripts

import (
	"os/exec"
	"bytes"
)

func Status() (*exec.Cmd, error) {
	c := exec.Command("git","status")
	out := new(bytes.Buffer)
	c.Stdout = out
	err := c.Run()
	if err != nil {
		return nil, err
	}
	return c, nil
}
