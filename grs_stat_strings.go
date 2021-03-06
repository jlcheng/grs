// Code generated by "stringer -type Dirstat,Branchstat,Indexstat -output grs_stat_strings.go"; DO NOT EDIT.

package grs

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GRSDIR_INVALID-0]
	_ = x[GRSDIR_VALID-1]
}

const _Dirstat_name = "GRSDIR_INVALIDGRSDIR_VALID"

var _Dirstat_index = [...]uint8{0, 14, 26}

func (i Dirstat) String() string {
	if i < 0 || i >= Dirstat(len(_Dirstat_index)-1) {
		return "Dirstat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Dirstat_name[_Dirstat_index[i]:_Dirstat_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BRANCH_UNKNOWN-0]
	_ = x[BRANCH_UPTODATE-1]
	_ = x[BRANCH_AHEAD-2]
	_ = x[BRANCH_BEHIND-3]
	_ = x[BRANCH_DIVERGED-4]
	_ = x[BRANCH_UNTRACKED-5]
}

const _Branchstat_name = "BRANCH_UNKNOWNBRANCH_UPTODATEBRANCH_AHEADBRANCH_BEHINDBRANCH_DIVERGEDBRANCH_UNTRACKED"

var _Branchstat_index = [...]uint8{0, 14, 29, 41, 54, 69, 85}

func (i Branchstat) String() string {
	if i < 0 || i >= Branchstat(len(_Branchstat_index)-1) {
		return "Branchstat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Branchstat_name[_Branchstat_index[i]:_Branchstat_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INDEX_UNKNOWN-0]
	_ = x[INDEX_MODIFIED-1]
	_ = x[INDEX_UNMODIFIED-2]
}

const _Indexstat_name = "INDEX_UNKNOWNINDEX_MODIFIEDINDEX_UNMODIFIED"

var _Indexstat_index = [...]uint8{0, 13, 27, 43}

func (i Indexstat) String() string {
	if i < 0 || i >= Indexstat(len(_Indexstat_index)-1) {
		return "Indexstat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Indexstat_name[_Indexstat_index[i]:_Indexstat_index[i+1]]
}
