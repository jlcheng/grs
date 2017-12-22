package status

import (
	"jcheng/grs/grs"
	"errors"
	"fmt"
)

type RepoStatus int

const (
	UNKNOWN RepoStatus = iota
	AHEAD
	CONFLICT
	LATEST
)

type Statusboard struct {
	repos map[string]RepoStatus
}

func (s *Statusboard) Status(repo *grs.Repo) (RepoStatus, error) {
	var r RepoStatus
	var ok bool
	if r, ok = s.repos[repo.Path]; !ok {
		return UNKNOWN, errors.New(fmt.Sprintf("repo not found [%v]", repo.Path))
	}
	return r, nil
}

func (s *Statusboard) SetStatus(repo *grs.Repo, status RepoStatus) {
	s.repos[repo.Path] = status
}

