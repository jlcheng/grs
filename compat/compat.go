package compat

import (
	"os/exec"
	"strings"
	"os"
)

// BeforeCmd sets up OS-specific changes
func BeforeCmd(cmd *exec.Cmd) {
	// https://github.com/git-for-windows/git/issues/1220#issuecomment-359302449
	// cygwin version of git will strip braces during globbing. Should be configurable TODO:1
	if strings.HasPrefix(os.Getenv("OS"),"Windows") {
		cmd.Env = append(cmd.Env, "CYGWIN=noglob")
		cmd.Env = append(cmd.Env, "MYSYS=noglob")
	}
}
