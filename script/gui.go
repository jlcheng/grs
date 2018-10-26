package script

import (
	"fmt"
	"jcheng/grs/status"
)

type AnsiGUI struct {
	clr bool // if true, clears screen before each iteration
}

func NewGUI(clr bool) AnsiGUI {
	return AnsiGUI{
		clr: clr,
	}
}

func (gui *AnsiGUI) Run(repos []status.Repo) {
	// setup/clear screen
	if gui.clr {
		fmt.Print("\033[2J\033[H")
	}

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v.\n",
			repo.Path, colorB(repo.Branch), colorI(repo.Index))
	}
}

func colorI(s status.Indexstat) string {
	if s == status.INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorB(s status.Branchstat) string {
	if s == status.BRANCH_UPTODATE {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}
