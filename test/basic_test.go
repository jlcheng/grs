package test

import (
	"fmt"
	"jcheng/grs/script"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args[3:]
	switch args[0] {
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
		os.Exit(0)
	case "false":
		os.Exit(1)
	default:
		os.Exit(1)
	}
	fmt.Println(os.Args[3:])
	fmt.Println("hello world")
}

func helperCommand(s ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd = exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestEcho(t *testing.T) {
	cmd := helperCommand("echo", "hello world")
	out, err := cmd.Output()
	if err != nil {
		t.Errorf("echo: %v", err)
	}
	if g, e := string(out), "hello world\n"; g != e {
		t.Errorf("echo: want %q, got %q", e, g)
	}
}

func TestFail(t *testing.T) {
	cmd := helperCommand("false")
	out, err := cmd.Output()
	if err != nil {
		if s := fmt.Sprintf("%v", err); s != "exit status 1" {
			t.Errorf("false: want [exit status 1], got [%v]", err)
		}
	} else {
		t.Errorf("false: want exit status 1, got exit status 0 with: %v", string(out))
	}
}

func TestReposFromString(t *testing.T) {
	var r []script.Repo

	r = script.ReposFromString("")
	if r[0].Path != "" {
		t.Error("TestReposFromgString")
	}

	r = script.ReposFromString("foo")
	if r[0].Path != "foo" {
		t.Error("TestReposFromString")
	}

	path0 := "/foo bar/fib"
	path1 := "file://fizz/fuzz"
	r = script.ReposFromString(path0 + string(os.PathListSeparator) + path1)
	if r[0].Path != path0 && r[1].Path != path1 {
		t.Error("TestReposFromString")
	}
}
