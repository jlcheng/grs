package status

import (
	"jcheng/grs/grs"
	"errors"
	"fmt"
)

type RepoStatus int

const (
	UNKNOWN RepoStatus = iota  // Repo status cannot be determined
	INVALID // Repo is not a valid Git repo
	AHEAD // Repo is ahead of remote
	CONFLICT // Repo is in conflict with remote
	LATEST // Repo is up-to-date with remote
	)
var statusStrings [5]string = [5]string{
	"UNKNOWN",
	"INVALID",
	"AHEAD",
	"CONFLICT",
	"LATEST",
}

func (s RepoStatus) String() string {
	return statusStrings[s]
}

type Statusboard struct {
	repos map[string]RepoStatus
}

func (s *Statusboard) Status(repo *grs.Repo) (RepoStatus, error) {
	var r RepoStatus
	var exists bool
	if r, exists = s.repos[repo.Path]; !exists {
		return UNKNOWN, errors.New(fmt.Sprintf("repo not found [%v]", repo.Path))
	}
	return r, nil
}

func (s *Statusboard) SetStatus(repo *grs.Repo, status RepoStatus) {
	s.repos[repo.Path] = status
}

