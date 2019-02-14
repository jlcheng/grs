package script

import (
	"errors"
	"fmt"
	"strings"
)

// AutoRebase rebases the current branch one change at a time and stops at the first error it encounters. The intent
// is that the current branch will now be as similar to @{upstream} as possible.
func (s *Script) AutoRebase() error {
	ctx := s.ctx
	runner := s.ctx.CommandRunner

	//  1. Identify merge-base
	git := ctx.GitExec
	p := "@{upstream}"
	cmd := runner.Command(git, "merge-base", "HEAD", p).WithDir(s.repo.Path)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	mergeBase := strings.TrimSpace(string(bytes))

	//  2. Identify the graph of child commits from merge-base to HEAD
	cmd = runner.Command(git, "rev-list", p, "^"+mergeBase).WithDir(s.repo.Path)
	bytes, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	revlist := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	//  3. Rebase current branch against each child in lineage
	var rebaseErr error
	for i := len(revlist) - 1; i >= 0 && rebaseErr == nil; i-- {
		commit := revlist[i]
		cmd = runner.Command(git, "rebase", commit).WithDir(s.repo.Path)
		_, err1 := cmd.CombinedOutput()
		if err1 != nil {
			//  4. Stop when conflict is detected
			rebaseErr = err1
			cmd = runner.Command(git, "rebase", "--abort").WithDir(s.repo.Path)
			bytes2, err2 := cmd.CombinedOutput()
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s", err2, string(bytes2)))
			}
		}
	}

	return rebaseErr
}
