package script

import (
	"jcheng/grs/grs"
	"jcheng/grs/status"
	"time"
)

// Fetch runs `git fetch`.
func (s *Script) Fetch() {
	if s.err != nil || s.repo.Dir != status.DIR_VALID {
		return
	}

	dbRepo := s.ctx.DB().FindOrCreateRepo(s.repo.Path)
	now := time.Now().Unix()
	if dbRepo.FetchedSec > (now - int64(s.ctx.MinFetchSec)) {
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
	dbRepo.FetchedSec = now
}
