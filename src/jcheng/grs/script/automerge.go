package script

import (
	"jcheng/grs/core"
	"jcheng/grs/status"
)

// AutoFFMerge runs `git merge --ff-only...` when the branch is behind and unmodified
func (s *Script) AutoFFMerge() bool {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil ||
		repo.Dir != status.DIR_VALID ||
		repo.Branch != status.BRANCH_BEHIND ||
		repo.Index != status.INDEX_UNMODIFIED {
		return false
	}

	git := ctx.GetGitExec()

	command := ctx.CommandRunner.Command(git, "merge", "--ff-only", "@{upstream}")
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		grs.Debug("git merge failed: %v\n%v\n", err, string(out))
		return false
	}
	return true
}
