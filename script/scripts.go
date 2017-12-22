package script

import (
	"jcheng/grs/grs"
	"jcheng/grs/status"
)

func GetRepoStatus(repo grs.Repo) status.RepoStatus {
	return status.UNKNOWN
}

type Script func(grs.Repo) status.RepoStatus

