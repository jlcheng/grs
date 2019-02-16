package script

import (
	"bytes"
	"errors"
	"fmt"
	"jcheng/grs/shexec"
	"os"
	"os/exec"
	"strings"
)

// GitTestHelper simplifies setting up Git repos on the filesysytem. If any method calls results in an error, all
// further method calls will be a no-op.
// It applies the Scanner.Err() technique as mentioned in https://blog.golang.org/errors-are-values
type GitTestHelper struct {
	err       error
	errCause  string
	git       string
	runner    shexec.CommandRunner
	debugExec bool
}

func NewGitTestHelper() *GitTestHelper {
	_, debugExec := os.LookupEnv("GRS_TEST_EXEC_DEBUG")

	return &GitTestHelper{
		err:       nil,
		git:       ResolveGit(),
		runner:    &shexec.ExecRunner{},
		debugExec: debugExec,
	}
}

func (s *GitTestHelper) Git() string {
	return s.git
}

func (s *GitTestHelper) Err() error {
	return s.err
}

func (s *GitTestHelper) ErrCause() string {
	return s.errCause
}

func (s *GitTestHelper) ErrString() string {
	return fmt.Sprintf("%v\n\n%v", s.errCause, s.err)
}

func (s *GitTestHelper) CommandRunner() shexec.CommandRunner {
	return s.runner
}

func (s *GitTestHelper) Mkdir(subdir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Mkdir(subdir, 0755); err != nil {
		s.err = err
		s.errCause = "mkdir " + subdir
		return false
	}
	return true
}

func (s *GitTestHelper) Chdir(dir string) bool {
	if s.err != nil {
		return false
	}
	if err := os.Chdir(dir); err != nil {
		s.err = err
		s.errCause = "chdir " + dir
		return false
	}
	return true
}

func (s *GitTestHelper) Touch(file string) bool {
	if s.err != nil {
		return false
	}
	f, err := os.Create(file)
	if err != nil {
		s.err = err
		s.errCause = "touch " + file
		return false
	}
	if f != nil {
		if err := f.Close(); err != nil {
			return false
		}
	}
	return true
}

func (s *GitTestHelper) SetContents(file, contents string) (ok bool) {
	if s.err != nil {
		return false
	}
	f, err := os.Create(file)
	if err != nil {
		s.err = err
		s.errCause = fmt.Sprintf("opening %v for write", file)
		return false
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			s.err = err2
			s.errCause = fmt.Sprintf("closing %v after write", file)
			ok = false
		}
	}()
	_, err = f.WriteString(contents)
	if err != nil {
		s.err = err
		s.errCause = "writing to " + file
		return false
	}
	ok = true
	return ok
}

func (s *GitTestHelper) Exec(first string, arg ...string) bool {
	if s.err != nil {
		return false
	}
	cmd := s.runner.Command(first, arg...)
	bytes, err := cmd.CombinedOutput()
	if s.debugExec {
		fmt.Println(first + strings.Join(arg, " "))
		fmt.Println(string(bytes))
	}
	if err != nil {
		s.err = errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
		s.errCause = first + " " + strings.Join(arg, " ")
		return false
	}
	return true
}

func (s *GitTestHelper) Add(path string) bool {
	return s.Exec(s.git, "add", path)
}

func (s *GitTestHelper) Commit(msg string) bool {
	if s.err != nil {
		return false
	}
	git := s.git
	return s.Exec(git, "commit", "-m", msg)
}

func (s *GitTestHelper) TouchAndCommit(file, msg string) bool {
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

