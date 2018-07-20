package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
)

// GetIndexStatus sets the rstat.index property to modified if there are uncommited changes or if the index has been
// modified
func (s *Script) GetIndexStatus() {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil || repo.Dir != status.DIR_VALID {
		return
	}

	git := ctx.GetGitExec()

	repo.Index = status.INDEX_UNKNOWN
	command := ctx.CommandRunner.Command(git, "ls-files", "--exclude-standard", "-om")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("ls-files failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		repo.Index = status.INDEX_MODIFIED
		return
	}

	command = ctx.CommandRunner.Command(git, "diff-index", "HEAD")
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("diff-index failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		repo.Index = status.INDEX_MODIFIED
		return
	}

	repo.Index = status.INDEX_UNMODIFIED
}
