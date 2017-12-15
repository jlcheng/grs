package grs

import (
	"os/exec"
	"bytes"
	"os"
)

type Repo struct {
	Path string
}

type Result struct {
	delegate *exec.Cmd
	Stdout string
}



func Status(repo Repo) (*Result, error) {
	cmd := new(Result)
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

func Pwd(repo Repo) (*Result, error) {
	cmd := new(Result)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}
	cmd.delegate = exec.Command("pwd")
	cmd.delegate.Stdout = new(bytes.Buffer)
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func Rebase(repo Repo) (*Result, error) {
	cmd := new(Result)
	err := os.Chdir(repo.Path)
	if err != nil {
		return cmd, err
	}
	buf := new(bytes.Buffer)
	cmd.delegate = exec.Command("git","fetch")
	cmd.delegate.Stdout = buf
	err = cmd.delegate.Run()
	if err != nil {
		return cmd, err
	}

	cmd.delegate = exec.Command("git", "rebase", "origin/master")
	cmd.delegate.Stdout = buf
	err = cmd.delegate.Run()
	if err != nil {
		cmd.delegate = exec.Command("git", "rebase", "--abort")
		cmd.delegate.Stdout = buf
		err = cmd.delegate.Run()
		if err != nil {
			return cmd, err
		}
		Debug("Unable to rebase. Original HEAD restored.")
		return cmd, nil
	}
	return cmd, nil
}

// Cmd takes a Repo to act on and returns the result of the command
type Cmd func(Repo) (*Result, error )



func (cmd *Result) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}