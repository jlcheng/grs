package gitscripts

import (
	"os/exec"
	"bytes"
)

type GitCmd struct {
	delegate *exec.Cmd
}

func Status() (*GitCmd, error) {
	cmd := new(GitCmd)
	cmd.delegate = exec.Command("sawew","stastus")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err := cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}
