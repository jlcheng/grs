package status

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
func NewRStat() *RStat {
	return &RStat{
		Dir: DIR_INVALID,
		Branch: BRANCH_UNKNOWN,
		Index: INDEX_UNKNOWN,
	}
}

