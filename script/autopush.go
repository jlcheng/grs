package script

import (
	"fmt"
	"jcheng/grs/shexec"
	"time"
)

// AutoPush runs `git push...` when the branch is ahead. It also autocommits changes.
func (s *Script) AutoPush() bool {
	repo := s.repo
	if s.err != nil ||
		repo.Dir != DIR_VALID ||
		repo.Index == INDEX_UNKNOWN ||
		!repo.PushAllowed {
		return false
	}
	if !(repo.Branch == BRANCH_AHEAD || repo.Branch == BRANCH_UPTODATE) {
		return false
	}

	shexec.Debug("git auto-push ok: %v", repo)
	ctx := s.ctx
	git := ctx.GitExec
	commitMsg := AutoPushGenCommitMsg(NewStdClock())
		var out []byte
	var err error
	var command shexec.Command
	if repo.Index == INDEX_MODIFIED {
		command := ctx.CommandRunner.Command(git, "add", "-A")
		if out, err = command.CombinedOutput(); err != nil {
			shexec.Debug("git add failed. %v, %v", err, string(out))
			return false
		}

		command = ctx.CommandRunner.Command(git, "commit", "-m", commitMsg)
		if out, err = command.CombinedOutput(); err != nil {
			shexec.Debug("git commit failed. commit-msg=%v\nerr-msg:%v\nout:%v", commitMsg, err, string(out))
			return false
		}
		repo.Index = INDEX_UNMODIFIED
	}


	command = ctx.CommandRunner.Command(git, "push")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("git push failed. %v, %v", err, string(out))
		return false
	}
	repo.Branch = BRANCH_UPTODATE
	return true
}

func AutoPushGenCommitMsg(clock Clock) string {
	return fmt.Sprintf("grs-autocommit:%v", clock.Now().Format(time.RFC3339))
}

// Clock interface allows one to mock functions of the time.Time type
type Clock interface {
	Now() time.Time
}
type StdClock struct {}
func (s *StdClock) Now() time.Time {
	return time.Now()
}
func NewStdClock() *StdClock {
	return &StdClock{}
}

type MockClock struct {
	NowRetval time.Time
}
func (s *MockClock) Now() time.Time {
	return s.NowRetval
}