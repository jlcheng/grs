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
		repo.Branch != BRANCH_AHEAD ||
		repo.Index == INDEX_UNKNOWN ||
		!repo.PushAllowed {
		return false
	}

	ctx := s.ctx
	git := ctx.GetGitExec()
	commitMsg := AutoPushGenCommitMsg(NewStdClock())
	command := ctx.CommandRunner.Command(git, "commit", "-m", commitMsg)
	var out []byte
	var err error
	if repo.Index == INDEX_MODIFIED {
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