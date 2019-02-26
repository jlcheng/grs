package script

import (
	"os"
)

type Script struct {
	ctx  *AppContext
	repo *Repo
	err  error
}

func NewScript(ctx *AppContext, repo *Repo) *Script {
	return &Script{ctx: ctx, repo: repo}
}

// BeforeScript sets up the Script object for future operations.
// It sets repo.Dir to DIR_VALID if the repo.Path exists and appears valid.
func (s *Script) BeforeScript() {
	if s.err != nil {
		return
	}

	if s.repo.Path == "" {
		s.repo.Dir = DIR_INVALID
		return
	}

	if finfo, err := os.Stat(s.repo.Path); err != nil || !finfo.IsDir() {
		s.repo.Dir = DIR_INVALID
		return
	}

	git := s.ctx.GitExec
	command := s.ctx.CommandRunner.Command(git, "show-ref", "-q", "--head", "HEAD").WithDir(s.repo.Path)
	if _, err := command.CombinedOutput(); err != nil {
		s.repo.Dir = DIR_INVALID
		return
	}
	s.repo.Dir = DIR_VALID
}

func (s *Script) GetRepo() *Repo {
	return s.repo
}

func (s *Script) Update() {
	s.GetCommitTime()
	s.GetIndexStatus()
	s.GetRepoStatus()
	s.BeforeScript()
}