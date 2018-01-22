package test

import (
	"testing"
	"os"
)

func TestWindowsFix(t *testing.T) {
	v, has_os := os.LookupEnv("OS")
	os.Setenv("OS", "Windows_NT")

	// see https://npf.io/2015/06/testing-exec-command/
	// TODO:2 runs this executable and call a function to verify the subprocess will have CYGWIN & MSSYS set

	if !has_os {
		os.Unsetenv("OS")
	} else {
		os.Setenv("OS", v)
	}
}

