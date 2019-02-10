package script

import "jcheng/grs/shexec"

// Fetch runs `git fetch`.
func (s *Script) Fetch() {
	if s.err != nil || s.repo.Dir != DIR_VALID {
		return
	}

	git := s.ctx.GitExec

	command := s.ctx.CommandRunner.Command(git, "fetch")
	if out, err := command.CombinedOutput(); err != nil {
		// fetch may have failed for common reasons, such as not adding yourxk ssh key to the agent
		shexec.Debug("git fetch failed: %v\n%v", err, string(out))
		return
	}
	shexec.Debug("git fetch ok: %v", s.repo.Path)
}
