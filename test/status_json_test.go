package test

import (
	"jcheng/grs/grsdb"
	"jcheng/grs/status"
	"testing"
)

func TestUpdateRStat(t *testing.T) {
	src := status.RStat{
		Branch: status.BRANCH_DIVERGED,
		Dir:    status.DIR_VALID,
		Index:  status.INDEX_MODIFIED,
	}
	var dest grsdb.RStat_Json
	dest.Update(src)
	if dest.Branch != src.Branch ||
		dest.Dir != src.Dir ||
		dest.Index != src.Index {
		t.Fatal("TestFromRStat src and dest do not match")
	}
}
