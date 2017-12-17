package grs

import (
	"testing"
	"os"
	"fmt"
	"os/exec"
)

func TestRepoPath(t *testing.T) {
	repo := Repo{"/foo/bar"}
	result, _ := RepoPath(repo)
	if result.String() != "/foo/bar" {
		t.Fail()
	}
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	fmt.Println("hello world")
}

func helperCommand(s ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd = exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelloWorld(t *testing.T) {
	cmd := helperCommand("echo", "hello world")
	out, err := cmd.Output()
	if err != nil {
		t.Errorf("echo: %v", err)
	}
	if g, e := string(out), "hello world\n"; g != e {
		t.Errorf("echo: want %q, got %q", e, g)
	}
}