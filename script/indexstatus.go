package script

import "jcheng/grs/shexec"

// GetIndexStatus sets the rstat.index property to modified if there are uncommited changes or if the index has been
// modified
func (s *Script) GetIndexStatus() {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil || repo.Dir != DIR_VALID {
		return
	}

	git := ctx.GitExec

	repo.Index = INDEX_UNKNOWN
	command := ctx.CommandRunner.Command(git, "ls-files", "--exclude-standard", "-om")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("ls-files failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		repo.Index = INDEX_MODIFIED
		return
	}

	command = ctx.CommandRunner.Command(git, "diff-index", "HEAD")
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("diff-index failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		repo.Index = INDEX_MODIFIED
		return
	}

	repo.Index = INDEX_UNMODIFIED
}
