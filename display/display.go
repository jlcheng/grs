package display

import (
	"jcheng/grs/status"
)

type Display interface {
	SummarizeRepos(repos []RepoStatus)
	Update()
}

type RepoStatus struct {
	Path     string
	Rstat    status.RStat
	Merged   bool
	MergeCnt int
}

type RepoStatusRoot struct {
	Statuses map[string]*RepoStatus
}

func NewRepoStatusRoot() *RepoStatusRoot {
	return &RepoStatusRoot{
		Statuses: make(map[string]*RepoStatus),
	}
}
