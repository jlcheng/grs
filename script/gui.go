package script

import (
	"fmt"
)

type AnsiGUI struct {
	clr bool // if true, clears screen before each iteration
}

func NewGUI(clr bool) AnsiGUI {
	return AnsiGUI{
		clr: clr,
	}
}

func (gui *AnsiGUI) Run(repos []Repo) {
	// setup/clear screen
	if gui.clr {
		fmt.Print("\033[2J\033[H")
	}

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v, %v.\n",
			repo.Path, colorB(repo.Branch), colorI(repo.Index), repo.CommitTime)
	}
}

func colorI(s Indexstat) string {
	if s == INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorB(s Branchstat) string {
	if s == BRANCH_UPTODATE {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}
