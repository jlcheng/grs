package grs

//go:generate stringer -type Dirstat,Branchstat,Indexstat -output grs_stat_strings.go

// Branchstat describes a repo's checked-out branch
type Branchstat int

const (
	BRANCH_UNKNOWN Branchstat = iota
	BRANCH_UPTODATE
	BRANCH_AHEAD
	BRANCH_BEHIND
	BRANCH_DIVERGED
	BRANCH_UNTRACKED
)

// Dirstat describes a repo's directory path
type Dirstat int

const (
	GRSDIR_INVALID Dirstat = iota
	GRSDIR_VALID
)

// Indexstat describes a repo's Git index
type Indexstat int

const (
	INDEX_UNKNOWN Indexstat = iota
	INDEX_MODIFIED
	INDEX_UNMODIFIED
)

// A GrsStat describes Git repository
type GrsStats struct {
	Branch     Branchstat // Whether the local repo has diverged
	Dir        Dirstat    // Whether the repo directory is valid
	Index      Indexstat  // Whehter the index of the local repo was modified
	CommitTime string     // A humand-readable string describing the time of the last commit
}

type GrsStatsOpt func(stats *GrsStats)

func NewGrsStats(options ...GrsStatsOpt) GrsStats {
	stats := &GrsStats{}
	for _, option := range options {
		option(stats)
	}
	return *stats
}

func WithBranchstat(branchstat Branchstat) GrsStatsOpt {
	return func(stats *GrsStats) {
		stats.Branch = branchstat
	}
}

func WithDirstat(dirstat Dirstat) GrsStatsOpt {
	return func(stats *GrsStats) {
		stats.Dir = dirstat
	}
}

func WithIndexstat(indexstat Indexstat) GrsStatsOpt {
	return func(stats *GrsStats) {
		stats.Index = indexstat
	}
}
