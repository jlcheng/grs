package display

import "jcheng/grs/status"

type Display interface {
	SummarizeRepos(repos []RepoStatus)
	Update()
}

type RepoStatus struct {
	Path string
	Rstat status.RStat
	Merged bool
}