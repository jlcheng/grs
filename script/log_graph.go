package script

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// LogGraph is a map of commit message (child) -> commit messages (list of parents)
type LogGraph map[string][]string

// String formats LogGraph as a JSON
func (lg LogGraph) String() string {
	b, err := json.MarshalIndent(lg, "", "  ")
	if err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", map[string][]string(lg))
}

// LogGraph creates a LogGraph object from the current working directory
func (s *GitTestHelper) LogGraph() (LogGraph, error) {
	var bytes []byte
	var err error
	command := s.runner.Command(s.git, "log", "--pretty=%h %s").WithDir(s.wd)
	if bytes, err = command.CombinedOutput(); err != nil {
		return nil, err
	}
	// id_msg is a map of commit_sha -> commit_message
	id_msg := make(map[string]string)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) < 1 {
			continue
		}
		b := strings.IndexRune(line, ' ')
		if b < 0 {
			return nil, fmt.Errorf("cannot parse: %s", line)
		}
		id_msg[strings.TrimSpace(line[:b])] = strings.TrimSpace(line[b:])
	}

	command = s.runner.Command(s.git, "log", "--pretty=%h %p").WithDir(s.wd)
	if bytes, err = command.CombinedOutput(); err != nil {
		return nil, err
	}
	lg := make(map[string][]string)
	for _, line := range strings.Split(string(bytes), "\n") {
		line := strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		// creates two arrays with the same size
		// ids is an array of [child_sha, parent_sha_1, parent_sha_2, ...]
		// msgs is an array of [child_commit_msg, parent_commit_msg_1, parent_commit_msg_2, ...]
		// then we can map child_msg_commit_msg to parent_commit_msgs
		ids := strings.Split(line, " ")
		if len(ids) < 1 {
			return nil, fmt.Errorf("cannot parse: %s", line)
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
