package script

import (
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"os"
)

type Script struct {
	ctx  *grs.AppContext
	repo *status.Repo
	err  error
}

func NewScript(ctx *grs.AppContext, repo *status.Repo) *Script {
	return &Script{ctx: ctx, repo: repo}
}

// BeforeScript sets up the Script object for future operations.
// First, it os.Chdir to the repo directory and validates the repo.
// Second, it sets rstat.Dir to `DIR_VALID` if a git command can be executed
func (s *Script) BeforeScript() {
	if s.err != nil {
		return
	}
	if err := os.Chdir(s.repo.Path); err != nil {
		s.repo.Dir = status.DIR_INVALID
		return
	}
	git := s.ctx.GetGitExec()
	command := s.ctx.CommandRunner.Command(git, "show-ref", "-q", "HEAD")
	if _, err := command.CombinedOutput(); err != nil {
		s.repo.Dir = status.DIR_INVALID
		return
	}
	s.repo.Dir = status.DIR_VALID
}
