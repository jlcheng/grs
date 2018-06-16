package display

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type AnsiDisplay struct {
	daemon bool
	writer *bufio.Writer
}

func NewAnsiDisplay(daemon bool, writer io.Writer) *AnsiDisplay {
	return &AnsiDisplay{daemon: daemon, writer: bufio.NewWriter(writer)}
}

func (ansi *AnsiDisplay) SummarizeRepos(repos []RepoVO) {
	// setup/clear screen
	if ansi.daemon {
		ansi.writer.WriteString("\033[2J\033[H")
	}

	// write out the status of each repository
	for _, repo := range repos {
		if repo.MergedSec == 0 {
			ansi.writer.WriteString(fmt.Sprintf("repo [%v] status is %v, %v.\n",
				repo.Repo.Path, repo.Repo.Branch, repo.Repo.Index))
		} else {
			ansi.writer.WriteString(fmt.Sprintf("repo [%v] status is %v, %v. Last merge on %v.\n",
				repo.Repo.Path, repo.Repo.Branch, repo.Repo.Index, fmtTime(repo.MergedSec)))
		}
	}

}

func fmtTime(sec int64) string {
	if sec == 0 {
		return "unknown"
	}
	return time.Unix(sec, 0).Format("Jan 2 15:04 PST")
}

func (ansi *AnsiDisplay) Update() {
	ansi.writer.Flush()
}
