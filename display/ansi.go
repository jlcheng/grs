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
		if repo.Merged {
			ansi.writer.WriteString(fmt.Sprintf("repo [%v] auto fast-foward to latest\n", repo.Path))
		} else {
			ansi.writer.WriteString(fmt.Sprintf("repo [%v] status is %v, %v\n", repo.Path, repo.Rstat.Branch,
				repo.Rstat.Index))
		}
	}

}

func (ansi *AnsiDisplay) Update() {
	ansi.writer.Flush()
}