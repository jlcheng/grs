package script

// AutoPush runs `git push...` when the branch is ahead and unmofified
func (s *Script) AutoPush() bool {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil ||
		repo.Dir != DIR_VALID ||
		repo.Branch != BRANCH_AHEAD ||
		repo.Index != INDEX_UNMODIFIED {
		return false
	}

	_ = ctx.GetGitExec()

	return false
}

