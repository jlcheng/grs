package script

import (
	"fmt"
	"jcheng/grs/shexec"
	"sort"
	"strings"
)

func (s *Script) LogGraph() (map[string][]string, error) {
	var bytes []byte
	var err error
	var command shexec.Command
	var lines []string
	git := s.ctx.GitExec

	command = s.ctx.CommandRunner.Command(git, "log", "--pretty=%h %s").WithDir(s.repo.Path)
	if bytes, err = command.CombinedOutput(); err != nil {
		s.repo.Dir = DIR_INVALID
		return nil, err
	}
	lines = strings.Split(string(bytes), "\n")
	id_msg := make(map[string]string)
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}
		b := strings.IndexRune(line, ' ')
		if b < 0 {
			return nil, fmt.Errorf("cannot parse: %s", string(bytes))
		}
		id := strings.TrimSpace(line[:b])
		msg := strings.TrimSpace(line[b:])
		id_msg[id] = msg
	}

	command = s.ctx.CommandRunner.Command(git, "log", "--pretty=%h %p").WithDir(s.repo.Path)
	if bytes, err = command.CombinedOutput(); err != nil {
		s.repo.Dir = DIR_INVALID
		return nil, err
	}
	lines = strings.Split(string(bytes), "\n")
	lg := make(map[string][]string)
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) <1 {
			continue
		}
		ids := strings.Split(line, " ")
		if len(ids) < 1 {
			return nil, fmt.Errorf("cannot parse: %s", string(bytes))
		}
		msgs := make([]string, len(ids))
		for idx, id := range ids {
			msg, ok := id_msg[id]
			if ok {
				msgs[idx] = msg
			}
		}
		parents := msgs[1:]
		sort.Slice(parents, func(i, j int) bool {
			return parents[i] < parents[j]
		})

		lg[msgs[0]] = parents
	}
	return lg, nil
}

