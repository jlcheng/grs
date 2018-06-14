package status

type Dirstat int

const (
	DIR_INVALID Dirstat = iota
	DIR_VALID
)

var dirstatStr = [DIR_VALID + 1]string{
	"INVALID",
	"VALID",
}

func (i Dirstat) String() string { return dirstatStr[i] }

type Branchstat int

const (
	BRANCH_UNKNOWN Branchstat = iota
	BRANCH_UPTODATE
	BRANCH_AHEAD
	BRANCH_BEHIND
	BRANCH_DIVERGED
	BRANCH_UNTRACKED
)

var branchstatdir = [BRANCH_UNTRACKED + 1]string{
	"UNKNOWN",
	"UP-TO-DATE",
	"AHEAD",
	"BEHIND",
	"DIVERGED",
	"UNTRACKED",
}

func (i Branchstat) String() string { return branchstatdir[i] }

type Indexstat int

const (
	INDEX_UNKNOWN Indexstat = iota
	INDEX_MODIFIED
	INDEX_UNMODIFIED
)

var indexstatdir = [INDEX_UNMODIFIED + 1]string{
	"UNKNOWN",
	"MODIFIED",
	"UNMODIFIED",
}

func (i Indexstat) String() string { return indexstatdir[i] }
