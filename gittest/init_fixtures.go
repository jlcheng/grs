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

	git := tctx.git
	if err = tctx.do(git, "init"); err != nil {
		return err
	}
	if err = tctx.touchAndCommit(".gitignore", "Commit_A"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("b.txt", "Commit_B"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "-b", "branch_A", "master~1"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("c.txt", "Commit_C"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "master"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("d.txt", "Commit_D"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "branch_A"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("e.txt", "Commit_E"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "-b", "branch_B", "master~1"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("f.txt", "Commit_F"); err != nil {
		return err
	}

	if err = tctx.touchAndCommit("g.txt", "Commit_G"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "master"); err != nil {
		return err
	}

	if err = tctx.do(git, "merge", "branch_B", "-m", "merge branch_B onto master"); err != nil {
		return err
	}

	if err = tctx.do(git, "checkout", "branch_A"); err != nil {
		return err
	}

	return err
}

func (tctx TestContext) do(first string, arg ...string) error {
	cmd := tctx.runner.Command(first, arg...)
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	} else if tctx.debugExec {
		fmt.Println(string(bytes))
	}

	return nil
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

func (tctx TestContext) touchAndCommit(file string, commitMsg string) (err error) {
	git := tctx.git
	if err = touch(file); err != nil {
		return err
	}
	if err = tctx.do(git, "add", "-A"); err != nil {
		return err
	}
	if err = tctx.do(git, "commit", "-m", commitMsg); err != nil {
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
