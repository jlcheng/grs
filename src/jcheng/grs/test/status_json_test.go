package test

import (
	"jcheng/grs/grsdb"
	"jcheng/grs/status"
	"testing"
)

func TestUpdateRStat(t *testing.T) {
	src := status.NewRepo("")
	src.Branch = status.BRANCH_DIVERGED
	src.Dir = status.DIR_VALID
	src.Index = status.INDEX_MODIFIED
	var dest grsdb.RStat_Json
	dest.Update(src)
	if dest.Branch != src.Branch ||
		dest.Dir != src.Dir ||
		dest.Index != src.Index {
		t.Fatal("TestFromRStat src and dest do not match")
	}
}
