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
	BEHIND // Repo is behind remote
	AHEAD // Repo is ahead of remote
	DIVERGED // Repo and remote have diverged - conflict unknown
	CONFLICT // Repo and remote have diverged - known conflict
	LATEST // Repo is up-to-date with remote
	)
var statusStrings [LATEST+1]string = [LATEST+1]string{
	"UNKNOWN",
	"INVALID",
	"BEHIND",
	"AHEAD",
	"DIVERGED",
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

