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

type Dirstat int
const (
	DIR_VALID Dirstat = iota
	DIR_INVALID
)
var dirstatStr [DIR_INVALID+1]string = [DIR_INVALID+1]string{
	"VALID",
	"INVALID",
}
func (i Dirstat) String() string { return dirstatStr[i] }
type Branchstat int
const (
	BRANCH_UNKNOWN Branchstat = iota
	BRANCH_UPTODATE
	BRANCH_AHEAD
	BRANCH_BEHIND
	BRANCH_DIVERGED
)
var branchstatdir [BRANCH_DIVERGED+1]string = [BRANCH_DIVERGED+1]string{
	"UNKNOWN",
	"UP-TO-DATE",
	"AHEAD",
	"BEHIND",
	"DIVERGED",
}
func (i Branchstat) String() string { return branchstatdir[i] }
type Indexstat int
const(
	INDEX_UNKNOWN Indexstat = iota
	INDEX_MODIFIED
	INDEX_UNMODIFIED
)
var indexstatdir [INDEX_UNMODIFIED+1]string = [INDEX_UNMODIFIED+1]string{
	"UNKNOWN",
	"MODIFIED",
	"UNMODIFIED",
}
func (i Indexstat) String() string { return indexstatdir[i] }
type RStat struct {
	Dir Dirstat
	Branch Branchstat
	Index Indexstat
}
func NewRStat() RStat {
	return RStat{
		Dir: DIR_INVALID,
		Branch: BRANCH_UNKNOWN,
		Index: INDEX_UNKNOWN,
	}
}

func (s RepoStatus) String() string {
	return statusStrings[s]
}

type entry struct {
	repo grs.Repo
	status RepoStatus
}
type Statusboard struct {
	repos map[string]entry
}


func (s *Statusboard) Status(path string) (RepoStatus, error) {
	var r entry
	var exists bool
	if r, exists = s.repos[path]; !exists {
		return UNKNOWN, errors.New(fmt.Sprintf("repo not found [%v]", path))
	}
	return r.status, nil
}

func (s *Statusboard) SetStatus(path string, status RepoStatus) {
	var r entry
	var exists bool
	if r, exists = s.repos[path]; !exists {
		s.repos[path] = entry{repo:grs.Repo{Path:path}, status:status}
	} else {
		r.status = status
	}
}

func (s *Statusboard) Repos() []string {
	var keys []string
	for k := range s.repos {
		keys = append(keys, k)
	}
	return keys
}

func NewStatusboard(repos ...grs.Repo) Statusboard {
	var s = Statusboard{}
	s.repos = make(map[string]entry, 0)
	for _, repo := range repos {
		s.SetStatus(repo.Path, UNKNOWN)
	}
	return s
}