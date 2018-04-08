package display

import (
	"bufio"
	"fmt"
	"io"
)

type AnsiDisplay struct {
	writer *bufio.Writer
}

func NewAnsiDisplay(writer io.Writer) (*AnsiDisplay) {
	return &AnsiDisplay{bufio.NewWriter(writer)}
}

func (ansi *AnsiDisplay) SummarizeRepos(repos []RepoStatus) {
	// setup/clear screen
	ansi.writer.WriteString("\033[2J\033[H")

	// write out the status of each repository
	for _, repo := range repos {
		ansi.writer.WriteString(fmt.Sprintf("repo [%v] status is %v, %v. %v auto-merges performed.\n",
			repo.Path, repo.Rstat.Branch, repo.Rstat.Index, repo.MergeCnt))
	}

}

func (ansi *AnsiDisplay) Update() {
	ansi.writer.Flush()
}