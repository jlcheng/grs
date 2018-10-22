package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
)

// Fetch runs `git fetch`.
func (s *Script) Fetch() {
	if s.err != nil || s.repo.Dir != status.DIR_VALID {
		return
	}

	git := s.ctx.GetGitExec()

	command := s.ctx.CommandRunner.Command(git, "fetch")
	if out, err := command.CombinedOutput(); err != nil {
		// fetch may have failed for common reasons, such as not adding yourxk ssh key to the agent
		grs.Debug("git fetch failed: %v\n%v", err, string(out))
		return
	}
	grs.Debug("git fetch ok: %v", s.repo.Path)
}
