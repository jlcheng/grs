package script

import (
	"bytes"
	"errors"
	"fmt"
	"jcheng/grs/shexec"
	"os"
	"os/exec"
)

// An ErrorExecRunner is a stateful utility that refuses to execute further Commands once an error occurs.
// It applies the Scanner.Err() technique as mentioned in https://blog.golang.org/errors-are-values
type ErrorExecRunner struct {
	err       error
	git       string
	runner    *shexec.ExecRunner
	debugExec bool
}

func NewExecRunner() *ErrorExecRunner {
	_, debugExec := os.LookupEnv("GRS_TEST_EXEC_DEBUG")

	return &ErrorExecRunner{
		err:       nil,
		git:       ResolveGit(),
		runner:    &shexec.ExecRunner{},
		debugExec: debugExec,
	}
}

func (s *ErrorExecRunner) Git() string {
	return s.git
}

func (s *ErrorExecRunner) Err() error {
	return s.err
}

func (s *ErrorExecRunner) ExecRunner() *shexec.ExecRunner {
	return s.runner
}

func (s *ErrorExecRunner) Mkdir(subdir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Mkdir(subdir, 0755); err != nil {
		s.err = err
		return false
	}
	return true
}

func (s *ErrorExecRunner) Chdir(dir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Chdir(dir); err != nil {
		s.err = err
		return false
	}
	return true
}

func (s *ErrorExecRunner) Touch(file string) bool {
	if s.err != nil {
		return false
	}
	f, err := os.Create(file)
	if err != nil {
		s.err = err
		return false
	}
	if f != nil {
		if err := f.Close(); err != nil {
			return false
		}
	}
	return true
}

func (s *ErrorExecRunner) SetContents(file, contents string) (ok bool) {
	if s.err != nil {
		return false
	}
	f, err := os.Create(file)
	if err != nil {
		s.err = err
		return false
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			s.err = err2
			ok = false
		}
	}()
	_, err = f.WriteString(contents)
	if err != nil {
		s.err = err
		return false
	}
	ok = true
	return ok
}

func (s *ErrorExecRunner) Exec(first string, arg ...string) bool {
	if s.err != nil {
		return false
	}
	cmd := s.runner.Command(first, arg...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		s.err = errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
		return false
	} else if s.debugExec {
		fmt.Println(string(bytes))
	}
	return true
}

func (s *ErrorExecRunner) Add(path string) bool {
	return s.Exec(s.git, "add", path)
}

func (s *ErrorExecRunner) Commit(msg string) bool {
	if s.err != nil {
		return false
	}
	git := s.git
	return s.Exec(git, "commit", "-m", msg)
}

func (s *ErrorExecRunner) TouchAndCommit(file, msg string) bool {
	return s.Touch(file) &&
		s.Add(file) &&
		s.Commit(msg)
}

func ResolveGit() string {
	git := "git"
	if tmp, ok := os.LookupEnv("GRS_TEST_GIT"); ok {
		git = tmp
	}
	return git
}

type Result struct {
	delegate *exec.Cmd
	Stdout   string
}

func (cmd *Result) String() string {
	return cmd.delegate.Stdout.(*bytes.Buffer).String()
}

