package script

type Repo struct {
	Path   string
	Dir    Dirstat
	Branch Branchstat
	Index  Indexstat
}

func NewRepo(path string) *Repo {
	return &Repo{
		Path:   path,
		Dir:    DIR_INVALID,
		Branch: BRANCH_UNKNOWN,
		Index:  INDEX_UNKNOWN,
	}
}
