package script

import "jcheng/grs/shexec"

// AutoFFMerge runs `git merge --ff-only...` when the branch is behind and unmodified
func (s *Script) AutoFFMerge() bool {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil ||
		repo.Dir != DIR_VALID ||
		repo.Branch != BRANCH_BEHIND ||
		repo.Index != INDEX_UNMODIFIED {
		return false
	}

	git := ctx.GitExec

	command := ctx.CommandRunner.Command(git, "merge", "--ff-only", "@{upstream}").WithDir(s.repo.Path)
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		shexec.Debug("git merge failed: %v\n%v\n", err, string(out))
		return false
	}
	return true
}
