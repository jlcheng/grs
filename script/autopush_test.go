package script

import (
	"testing"
)

// TODO JCHENG
func TestAutoPush_Noop(t *testing.T) {
	verify_AutoPush_NoGitExec(t, DIR_INVALID, BRANCH_BEHIND, INDEX_UNMODIFIED)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_AHEAD, INDEX_UNMODIFIED)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_UNKNOWN, INDEX_UNMODIFIED)

	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_UNKNOWN)
	verify_AutoPush_NoGitExec(t, DIR_VALID, BRANCH_BEHIND, INDEX_MODIFIED)
}

// TODO JCHENG
func TestAutoPush_Ok(t *testing.T) {

}

// TODO JCHENG
func verify_AutoPush_NoGitExec(t *testing.T, dir Dirstat, branch Branchstat, index Indexstat) {

}
