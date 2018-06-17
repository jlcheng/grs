package script

import (
	"errors"
	"fmt"
	"strings"
)

func (s *Script) AutoRebase() error {
	ctx := s.ctx
	runner := s.ctx.CommandRunner

	//  2. Identify merge-base
	git := ctx.GetGitExec()
	p := "@{upstream}"
	cmd := runner.Command(git, "merge-base", "HEAD", p)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	mergeBase := strings.TrimSpace(string(bytes))

	//  3. Identify the graph of child commits from merge-base to HEAD
	cmd = runner.Command(git, "rev-list", p, "^"+mergeBase)
	bytes, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%v %v", err, string(bytes)))
	}
	revlist := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	//  5. Rebase current branch against each child in lineage
	for i := len(revlist) - 1; i >= 0; i-- {
		commit := revlist[i]
		cmd = runner.Command(git, "rebase", commit)
		_, err1 := cmd.CombinedOutput()
		if err1 != nil {
			cmd = runner.Command(git, "rebase", "--abort")
			bytes2, err2 := cmd.CombinedOutput()
			if err != nil {
				return errors.New(fmt.Sprintf("%s %s", err2, string(bytes2)))
			}
		}
	}

	//  6. Stop when conflict is detected
	return nil
}
