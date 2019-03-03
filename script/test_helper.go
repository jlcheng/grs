package script

import (
	"bytes"
	"errors"
	"fmt"
	"jcheng/grs/shexec"
	"os"
	"os/exec"
	"path"
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
	wd        string
}
type option func(*GitTestHelper)

func NewGitTestHelper(options ...option) *GitTestHelper {
	retval := &GitTestHelper{
		err: nil,
		git: ResolveGit(),
		runner: &shexec.ExecRunner{},
		debugExec: false,
	}
	for _, o := range options {
		o(retval)
	}

	if retval.wd == "" {
		retval.wd, _ = os.Getwd()
	}
	if retval.wd == "" {
		retval.wd = os.TempDir()
	}

	return retval
}

func WithDebug(debugExec bool) option {
	return func(g *GitTestHelper) {
		g.debugExec = debugExec
	}
}

func WithWd(wd string) option {
	return func(g *GitTestHelper) {
		g.wd = wd
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

func (s *GitTestHelper) Mkdir(dir string) bool {
	if s.err != nil {
		return false
	}
	target := s.toAbsPath(dir)
	if err := os.Mkdir(target, 0755); err != nil {
		s.err = err
		s.errCause = "mkdir " + target
		return false
	}
	return true
}

func (s *GitTestHelper) Chdir(dir string) bool {
	if s.err != nil {
		return false
	}
	target := s.toAbsPath(dir)
	s.wd = target
	return true
}

func (s *GitTestHelper) Getwd() string {
	return s.wd
}

func (s *GitTestHelper) Touch(file string) bool {
	if s.err != nil {
		return false
	}
	target := s.toAbsPath(file)
	f, err := os.Create(target)
	if err != nil {
		s.err = err
		s.errCause = "touch " + target
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
	target := s.toAbsPath(file)
	f, err := os.Create(target)
	if err != nil {
		s.err = err
		s.errCause = fmt.Sprintf("opening %v for write", target)
		return false
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			s.err = err2
			s.errCause = fmt.Sprintf("closing %v after write", target)
			ok = false
		}
	}()
	_, err = f.WriteString(contents)
	if err != nil {
		s.err = err
		s.errCause = "writing to " + target
		return false
	}
	ok = true
	return ok
}

// Exec executes the given command using GetWd() as the working directory.
// For example, `NewGitHelper(WithWd("/tmp")); Exec("ls")` will list the contents of `/tmp`.
func (s *GitTestHelper) Exec(first string, arg ...string) bool {
	if s.err != nil {
		return false
	}
	cmd := s.runner.Command(first, arg...).WithDir(s.Getwd())
	bytes, err := cmd.CombinedOutput()
	if s.debugExec {
		words := append([]string{">>>", first}, arg...)
		fmt.Println(strings.Join(words, " "))
		fmt.Println(string(bytes))
	}
	if err != nil {
		s.err = errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
		s.errCause = first + " " + strings.Join(arg, " ")
		return false
	}
	return true
}

// RunGit is a convenience method for calling "git ..."
func (s *GitTestHelper) RunGit(args ...string) bool {
	return s.Exec(s.Git(), args...)
}

// NewRepoPair creates two repos under the specified directory named source and dest. The source repo is initialized
// as a bare repository. The dest repo is used to add an initial commit and push said commit to source. This method
// changes the helper's working directory to the "dest" directory if possible.
func (s *GitTestHelper) NewRepoPair(basedir string) {
	git := s.Git()
	target := s.toAbsPath(basedir)
	s.Chdir(target)
	s.Mkdir("source")
	s.Chdir("source")
	s.Exec(git, "init", "--bare")
	s.Chdir("..")
	s.Exec(git, "clone", "source", "dest")

	s.Chdir("dest")
	s.TouchAndCommit("init.txt", "init")
	s.Exec(git, "push", "origin")
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

// Given the current working directory and a target path, turn the target path to an absolute path
func (s *GitTestHelper) toAbsPath(name string) string {
	if path.IsAbs(name) {
		return name
	}
	wd := s.Getwd()
	if wd == "" {
		wd, _ = os.Getwd()
		if wd == "" {
			wd = os.TempDir()
		}
	}

	return path.Join(s.wd, name)
}