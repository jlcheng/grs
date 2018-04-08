package display

import (
	"jcheng/grs/status"
)

type Display interface {
	SummarizeRepos(repos []RepoVO)
	Update()
}

type RepoVO struct {
	Path      string
	Rstat     status.RStat
	Merged    bool
	MergeCnt  int
	MergedSec int64
}

type RepoListVO struct {
	Statuses map[string]*RepoVO
}
