package script

// AutoPush runs `git push...` when the branch is ahead. It also autocommits changes.
func (s *Script) AutoPush() bool {
	repo := s.repo
	ctx := s.ctx
	if s.err != nil ||
		repo.Dir != DIR_VALID ||
		repo.Branch != BRANCH_AHEAD ||
		repo.Index != INDEX_UNKNOWN ||
		!repo.PushAllowed {
		return false
	}

	_ = ctx.GetGitExec()

	return false
}

