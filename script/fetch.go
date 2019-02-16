package script

import (
	"jcheng/grs/base"
)

// Fetch runs `git fetch`.
func (s *Script) Fetch() {
	if s.err != nil || s.repo.Dir != DIR_VALID {
		return
	}

	git := s.ctx.GitExec

	command := s.ctx.CommandRunner.Command(git, "fetch").WithDir(s.repo.Path)
	if out, err := command.CombinedOutput(); err != nil {
		// fetch may have failed for common reasons, such as not adding yourxk ssh key to the agent
		base.Debug("git fetch failed: %v\n%v", err, string(out))
		return
	}
	base.Debug("git fetch ok: %v", s.repo.Path)
}
