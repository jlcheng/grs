package gitscripts

import (
	"os/exec"
	"bytes"
	"os"
)

type Repo struct {
	Path string
}

type GitCmd struct {
	delegate *exec.Cmd
	Stdout string
}

func Status(repo Repo) (*GitCmd, error) {
	cmd := new(GitCmd)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}
	cmd.delegate = exec.Command("git","status")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func Pwd(repo Repo) (*GitCmd, error) {
	cmd := new(GitCmd)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}
	cmd.delegate = exec.Command("ls","-alth")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (cmd *GitCmd) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}