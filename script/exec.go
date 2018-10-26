package script

import (
	"errors"
	"fmt"
	"jcheng/grs/shexec"
	"os"
)

// An TestExecRunner is a stateful utility that refuses to execute further Commands once an error occurs.
// It applies the Scanner.Err() technique as mentioned in https://blog.golang.org/errors-are-values
type ExecRunner struct {
	err       error
	git       string
	runner    *shexec.ExecRunner
	debugExec bool
}

func NewExecRunner() *ExecRunner {
	_, debugExec := os.LookupEnv("GRS_TEST_EXEC_DEBUG")

	return &ExecRunner{
		err:       nil,
		git:       ResolveGit(),
		runner:    &shexec.ExecRunner{},
		debugExec: debugExec,
	}
}

func (s *ExecRunner) Git() string {
	return s.git
}

func (s *ExecRunner) Err() error {
	return s.err
}

func (s *ExecRunner) Runner() *shexec.ExecRunner {
	return s.runner
}

func (s *ExecRunner) Mkdir(subdir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Mkdir(subdir, 0755); err != nil {
		s.err = err
		return false
	}
	return true
}

func (s *ExecRunner) Chdir(dir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Chdir(dir); err != nil {
		s.err = err
		return false
	}
	return true
}

func (s *ExecRunner) Touch(file string) bool {
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

func (s *ExecRunner) SetContents(file, contents string) (ok bool) {
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

func (s *ExecRunner) Exec(first string, arg ...string) bool {
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

func (s *ExecRunner) Add(path string) bool {
	return s.Exec(s.git, "add", path)
}

func (s *ExecRunner) Commit(msg string) bool {
	if s.err != nil {
		return false
	}
	git := s.git
	return s.Exec(git, "commit", "-m", msg)
}

func (s *ExecRunner) TouchAndCommit(file, msg string) bool {
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
