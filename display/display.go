package display

import "jcheng/grs/status"

type Display interface {
	SummarizeRepos(repos []RepoVO)
	Update()
}

type RepoVO struct {
	Repo      status.Repo
	Merged    bool
	MergeCnt  int
	MergedSec int64
}

type RepoListVO struct {
	Statuses map[string]*RepoVO
}
