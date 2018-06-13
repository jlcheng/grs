package test

import (
	"fmt"
	"jcheng/grs/compat"
	"os"
	"os/exec"
	"testing"
)

func TestWindowsFix(t *testing.T) {
	v, has_os := os.LookupEnv("OS")
	os.Setenv("OS", "Windows_NT")

	cmd := helperExec(t, "TestWindowsFixHelper", "@{upstream}...HEAD")
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
	var arg string = ""
	for idx, elem := range os.Args {
		if elem == "--" && idx < len(os.Args)-1 {
			arg = os.Args[idx+1]
		}
	}
	if expected := "@\\{upstream\\}...HEAD"; arg != expected {
		t.Errorf("expected arg to be [%v] but got [%v]", expected, arg)
	}
}

func helperExec(t *testing.T, tname string, s ...string) *exec.Cmd {
	cs := []string{fmt.Sprintf("-test.run=%v", tname), "--"}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
