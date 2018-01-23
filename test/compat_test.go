package test

import (
	"testing"
	"os"
	"os/exec"
	"fmt"
	"jcheng/grs/compat"
)

func TestWindowsFix(t *testing.T) {
	v, has_os := os.LookupEnv("OS")
	os.Setenv("OS", "Windows_NT")

	cmd := helperExec(t, "TestWindowsFixHelper")
	compat.BeforeCmd(cmd)
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(b))
	}

	if !has_os {
		os.Unsetenv("OS")
	} else {
		os.Setenv("OS", v)
	}
}

func TestWindowsFixHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("CYGWIN") != "noglob" {
		t.Error("expected ENV[CYGWIN]=noglob")
	}
	if os.Getenv("MSYS") != "noglob" {
		t.Error("expected ENV[MSYS]=noglob")
	}
}

func helperExec(t *testing.T, tname string, s ...string) *exec.Cmd {
	cs := []string{fmt.Sprintf("-test.run=%v", tname), "--"}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}