package gittest

import (
	"jcheng/grs/grs"
	"os"
	"errors"
	"fmt"
)

type TestContext struct {
	git       string
	runner    *grs.ExecRunner
	debugExec bool
}

func NewTestContext() TestContext {
	_, debugExec := os.LookupEnv("GRS_TEST_EXEC_DEBUG")

	return TestContext{
		git:       ResolveGit(),
		runner:    &grs.ExecRunner{},
		debugExec: debugExec,
	}
}

func (s *TestContext) GetRunner() *grs.ExecRunner {
	return s.runner
}

func (tctx TestContext) TouchAndCommit(file string, commitMsg string) {
	git := tctx.git
	if err := Touch(file); err != nil {
		panic(err)
	}
	tctx.Exec(git, "add", file)
	tctx.Exec(git, "commit", "-m", commitMsg)
}

func (tctx TestContext) Mkdir(subdir string) {
	if err := os.Mkdir(subdir, 0755); err != nil {
		panic(err)
	}
}

func (tctx TestContext) Chdir(dir string) {
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func (tctx TestContext) Git() string {
	return tctx.git
}

func Touch(file string) error {
	f, err := os.Create(file)
	if f != nil {
		f.Close()
	} else if err != nil {
		return err
	}
	return nil
}

func ResolveGit() string {
	git := "git"
	if tmp, ok := os.LookupEnv("GRS_TEST_GIT"); ok {
		git = tmp
	}
	return git
}


func (tctx TestContext) Exec(first string, arg ...string) {
	cmd := tctx.runner.Command(first, arg...)
	if bytes, err := cmd.CombinedOutput(); err != nil {
		panic(errors.New(fmt.Sprintf("%v %v", err, string(bytes))))
	} else if tctx.debugExec {
		fmt.Println(string(bytes))
	}
}