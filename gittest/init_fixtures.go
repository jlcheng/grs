package gittest

import (
	"errors"
	"fmt"
	"jcheng/grs/grs"
	"os"
)

type TestContext struct {
	git       string
	runner    *grs.ExecRunner
	debugExec bool
}

func (s *TestContext) GetRunner() *grs.ExecRunner {
	return s.runner
}

// TestID:it_test_1 Sets up a git repository for it_test_1, rooted at tmpdir.
func InitTest1(tctx TestContext, tmpdir string) (err error) {
	if err := os.Chdir(tmpdir); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("error %v", r)
			}
		}
	}()

	git := tctx.git
	tctx.do(git, "init")
	tctx.touchAndCommit(".gitignore", "Commit_A")
	tctx.touchAndCommit("b.txt", "Commit_B")
	tctx.do(git, "checkout", "-b", "branch_A", "master~1")
	tctx.touchAndCommit("c.txt", "Commit_C")
	tctx.do(git, "checkout", "master")
	tctx.touchAndCommit("d.txt", "Commit_D")
	tctx.do(git, "checkout", "branch_A")
	tctx.touchAndCommit("e.txt", "Commit_E")
	tctx.do(git, "checkout", "-b", "branch_B", "master~1")
	tctx.touchAndCommit("f.txt", "Commit_F")
	tctx.touchAndCommit("g.txt", "Commit_G")
	tctx.do(git, "checkout", "master")
	tctx.do(git, "merge", "branch_B", "-m", "merge branch_B onto master")
	tctx.do(git, "checkout", "branch_A")
	return err
}

func (tctx TestContext) do(first string, arg ...string) {
	cmd := tctx.runner.Command(first, arg...)
	if bytes, err := cmd.CombinedOutput(); err != nil {
		panic(errors.New(fmt.Sprintf("%v %v", err, string(bytes))))
	} else if tctx.debugExec {
		fmt.Println(string(bytes))
	}
}

func (tctx TestContext) touchAndCommit(file string, commitMsg string) {
	git := tctx.git
	if err := touch(file); err != nil {
		panic(err)
	}
	tctx.do(git, "add", file)
	tctx.do(git, "commit", "-m", commitMsg)
}

func touch(file string) error {
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

func NewTestContext() TestContext {
	_, debugExec := os.LookupEnv("GRS_TEST_EXEC_DEBUG")

	return TestContext{
		git:       ResolveGit(),
		runner:    &grs.ExecRunner{},
		debugExec: debugExec,
	}
}
