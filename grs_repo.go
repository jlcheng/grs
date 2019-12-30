package grs

import (
	"fmt"
	"jcheng/grs/base"
	"jcheng/grs/shexec"
	"os"
	"strings"
)

// GrsRepo represents a git repository on your local file system.
type GrsRepo struct {
	stats         GrsStats
	err           error  // err is set when a method returns an error; may prevent further methods from running
	git           string // Path to Git executable
	local         string // Path to the local clone of the repo
	pushAllowed   bool   // If true, GrsRepo is allowed to push changes to remote
	commandRunner shexec.CommandRunner
}

// GrsRepoOpt provides functional options
type GrsRepoOpt func(gr *GrsRepo)

// WithCommandRunnerGrsRepo is an option for the CommandRunner
func WithCommandRunnerGrsRepo(commandRunner shexec.CommandRunner) GrsRepoOpt {
	return func(gr *GrsRepo) {
		gr.commandRunner = commandRunner
	}
}

// WithLocalGrsRepo is an option for the repo's path on the local file system
func WithLocalGrsRepo(local string) GrsRepoOpt {
	return func(gr *GrsRepo) {
		gr.local = local
	}
}

// WithPushAllowed is an option to enable auto-push
func WithPushAllowed(pushAllowed bool) GrsRepoOpt {
	return func(gr *GrsRepo) {
		gr.pushAllowed = pushAllowed
	}
}

// WithStats is an option to initialize a repo's stats
func WithStats(stats GrsStats) GrsRepoOpt {
	return func(gr *GrsRepo) {
		gr.stats = stats
	}
}

// WithError is an option to initiailize a repo's error
func WithError(error error) GrsRepoOpt {
	return func(gr *GrsRepo) {
		gr.err = error
	}
}

// NewGrsRepo returns an instance of GresRepo
func NewGrsRepo(options ...GrsRepoOpt) GrsRepo {
	retval := GrsRepo{
		git: "git",
	}
	for _, option := range options {
		option(&retval)
	}
	return retval
}

func (gr *GrsRepo) IsPushAllowed() bool {
	return gr.pushAllowed
}

// UpdateCommitTime reads the last commit time from Git
func (gr *GrsRepo) UpdateCommitTime() {
	if gr.err != nil || gr.stats.Dir != GRSDIR_VALID {
		return
	}

	var command shexec.Command
	var out []byte
	var err error
	var statsPtr = &gr.stats

	command = gr.commandRunner.Command(gr.git, "log", "-1", "--format=%cr").WithDir(gr.local)

	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, " failed: %v\n%v\n", err, string(out))
		gr.err = err
		statsPtr.CommitTime = "Unknown"
	}
	cmdStr := "+ git log -1 --format=%%cr"
	base.DebugFull("", gr.local, fmt.Sprintf("%s... %s", cmdStr, strings.TrimSpace(string(out))))
	statsPtr.CommitTime = strings.Trim(string(out), "\n")
}

// UpdateDirstat sets up the Script object for future operations.
// It sets repo.Dir to DIR_VALID if the repo.Path exists and appears valid.
func (gr *GrsRepo) UpdateDirstat() {
	if gr.err != nil {
		return
	}
	var statsPtr = &gr.stats
	if gr.local == "" {
		gr.err = fmt.Errorf("local not specified")
		statsPtr.Dir = GRSDIR_INVALID
		return
	}

	if finfo, err := os.Stat(gr.local); err != nil || !finfo.IsDir() {
		gr.err = fmt.Errorf("local not a directory")
		statsPtr.Dir = GRSDIR_INVALID
		return
	}

	git := gr.git
	command := gr.commandRunner.Command(git, "show-ref", "-q", "--head", "HEAD").WithDir(gr.local)
	if _, err := command.CombinedOutput(); err != nil {
		gr.err = err
		statsPtr.Dir = GRSDIR_INVALID
		return
	}
	statsPtr.Dir = GRSDIR_VALID
}

// UpdateRepoStatus update the "branch" status of a *GrsRepo
func (gr *GrsRepo) UpdateRepoStatus() {
	if gr.err != nil || gr.stats.Dir != GRSDIR_VALID {
		return
	}

	git := gr.git
	var command shexec.Command
	var out []byte
	var err error
	var statsPtr = &gr.stats

	command = gr.commandRunner.Command(git, "rev-parse", "@{upstream}").WithDir(gr.local)
	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, "UpdateRepoStatus: no upstream detected. %s, %s", err, strings.TrimSpace(string(out)))
		gr.err = err
		statsPtr.Branch = BRANCH_UNTRACKED
		return
	}

	cmdLine := "+ git rev-list --left-right --count @{upstream}...HEAD"
	command = gr.commandRunner.Command(git, "rev-list", "--left-right", "--count", "@{upstream}...HEAD").WithDir(gr.local)
	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, "git rev-list failed: %s %v\n%s", cmdLine, err, string(out))
		gr.err = err
		statsPtr.Dir = GRSDIR_INVALID
		statsPtr.Branch = BRANCH_UNKNOWN
		return
	}
	cmpOut := strings.TrimSpace(string(out))
	if cmpOut != "0\t0" {
		base.DebugFull("", gr.local, "%s... %s", cmdLine, cmpOut)
	}
	diff, err := parseRevList(out)
	if err != nil {
		base.DebugFull("", gr.local, "cannot parse `git rev-list...` output: %q", string(out))
		gr.err = err
		statsPtr.Dir = GRSDIR_INVALID
		statsPtr.Branch = BRANCH_UNKNOWN
		return
	}

	if diff.remote == 0 {
		if diff.local == 0 {
			statsPtr.Branch = BRANCH_UPTODATE
		} else {
			statsPtr.Branch = BRANCH_AHEAD
		}
	} else {
		if diff.local == 0 {
			statsPtr.Branch = BRANCH_BEHIND
		} else {
			statsPtr.Branch = BRANCH_DIVERGED
		}
	}
}

// UpdateIndexStatus updates the INDEX status of a GrsRepo
func (gr *GrsRepo) UpdateIndexStatus() {
	if gr.err != nil || gr.stats.Dir != GRSDIR_VALID {
		return
	}

	var statsPtr = &gr.stats

	statsPtr.Index = INDEX_UNKNOWN
	command := gr.commandRunner.Command(gr.git, "ls-files", "--exclude-standard", "-om").WithDir(gr.local)
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, "ls-files failed: %v\n%v\n", err, string(out))
		return
	}
	if len(out) != 0 {
		statsPtr.Index = INDEX_MODIFIED
		return
	}

	command = gr.commandRunner.Command(gr.git, "diff-index", "HEAD").WithDir(gr.local)
	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, "diff-index failed: %v\n%v\n", err, string(out))
		gr.err = err
		return
	}
	if len(out) != 0 {
		statsPtr.Index = INDEX_MODIFIED
		return
	}

	statsPtr.Index = INDEX_UNMODIFIED
}

// Fetch runs git fetch on the GrsRepo instance
func (gr *GrsRepo) Fetch() {
	if gr.err != nil || gr.stats.Dir != GRSDIR_VALID {
		return
	}

	command := gr.commandRunner.Command(gr.git, "fetch").WithDir(gr.local)
	if out, err := command.CombinedOutput(); err != nil {
		// fetch may have failed for common reasons, such as not adding your ssh key to the agent
		base.DebugFull("", gr.local, "git fetch failed: %v\n%v", err, string(out))
		gr.err = err
		return
	}
	base.DebugFull("", gr.local, "git fetch ok: %v", gr.local)
}

// AutoPush attempts to commit any changes and push them to the remote repo
func (gr *GrsRepo) AutoPush() {
	if gr.err != nil ||
		!gr.pushAllowed ||
		gr.stats.Dir != GRSDIR_VALID ||
		gr.stats.Index == INDEX_UNKNOWN ||
		gr.stats.Branch == BRANCH_UNKNOWN ||
		gr.stats.Branch == BRANCH_UNTRACKED {
		return
	}

	base.DebugFull("", gr.local, "AutoPush eligible")
	commitMsg := AutoPushGenCommitMsg(NewStdClock())
	var out []byte
	var err error
	var command shexec.Command
	if gr.stats.Index == INDEX_MODIFIED {
		command := gr.commandRunner.Command(gr.git, "add", "-A").WithDir(gr.local)
		if out, err = command.CombinedOutput(); err != nil {
			base.DebugFull("", gr.local, "git add failed. %v, %v", err, string(out))
			gr.err = err
			return
		}

		command = gr.commandRunner.Command(gr.git, "commit", "-m", commitMsg).WithDir(gr.local)
		if out, err = command.CombinedOutput(); err != nil {
			base.DebugFull("", gr.local, "git commit failed. commit-msg=%v\nerr-msg:%v\nout:%v", commitMsg, err, string(out))
			gr.err = err
			return
		}
	}

	if gr.stats.Branch == BRANCH_UPTODATE || gr.stats.Branch == BRANCH_AHEAD {
		command = gr.commandRunner.Command(gr.git, "push").WithDir(gr.local)
		if out, err = command.CombinedOutput(); err != nil {
			base.DebugFull("", gr.local, "git push failed. %v, %v", err, string(out))
			gr.err = err
			return
		}
		base.DebugFull("", gr.local, "AutoPush complete")
	}

	gr.UpdateIndexStatus()
	gr.UpdateRepoStatus()
}

// AutoRebase runs a smarter version of 'git --rebase'
func (gr *GrsRepo) AutoRebase() {
	base.DebugFull("", gr.local, "AutoRebase start")
	if gr.err != nil ||
		gr.stats.Dir != GRSDIR_VALID ||
		gr.stats.Index == INDEX_UNKNOWN ||
		gr.stats.Branch == BRANCH_UNKNOWN ||
		gr.stats.Branch == BRANCH_UNTRACKED {
		base.DebugFull("", gr.local, "AutoRebase aborted")
		return
	}

	//  1. Identify merge-base
	p := "@{upstream}"
	cmd := gr.commandRunner.Command(gr.git, "merge-base", "HEAD", p).WithDir(gr.local)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		gr.err = fmt.Errorf("%v %v", err, string(bytes))
		return
	}
	mergeBase := strings.TrimSpace(string(bytes))

	//  2. Identify the graph of child commits from merge-base to HEAD
	cmd = gr.commandRunner.Command(gr.git, "rev-list", p, "^"+mergeBase).WithDir(gr.local)
	bytes, err = cmd.CombinedOutput()
	if err != nil {
		gr.err = fmt.Errorf("%v %v", err, string(bytes))
		return
	}
	revlist := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	//  3. Rebase current branch against each child in lineage
	var rebaseErr error
	for i := len(revlist) - 1; i >= 0 && rebaseErr == nil; i-- {
		commit := revlist[i]
		cmd = gr.commandRunner.Command(gr.git, "rebase", commit).WithDir(gr.local)
		_, err1 := cmd.CombinedOutput()
		if err1 != nil {
			//  4. Stop when conflict is detected
			base.DebugFull("", gr.local, "%s:AutoRebase conflict at %s", gr.local, commit)
			rebaseErr = err1
			cmd = gr.commandRunner.Command(gr.git, "rebase", "--abort").WithDir(gr.local)
			bytes2, err2 := cmd.CombinedOutput()
			if err2 != nil {
				gr.err = fmt.Errorf("%v %v", err2, string(bytes2))
				return
			}
		}
	}
}

// AutoFFMerge runs git merge --ff-only
func (gr *GrsRepo) AutoFFMerge() {
	if gr.err != nil ||
		gr.stats.Dir == GRSDIR_INVALID ||
		gr.stats.Branch == BRANCH_AHEAD ||
		gr.stats.Branch == BRANCH_DIVERGED ||
		gr.stats.Branch == BRANCH_UNKNOWN ||
		gr.stats.Branch == BRANCH_UNTRACKED ||
		gr.stats.Index != INDEX_UNMODIFIED {
		return
	}

	command := gr.commandRunner.Command(gr.git, "merge", "--ff-only", "@{upstream}").WithDir(gr.local)
	var out []byte
	var err error
	if out, err = command.CombinedOutput(); err != nil {
		base.DebugFull("", gr.local, "git merge failed: %v\n%v\n", err, string(out))
		gr.err = err
	}
}

// Update updates the state of the GrsRepo
func (gr *GrsRepo) Update() {
	gr.UpdateDirstat()
	gr.UpdateCommitTime()
	gr.UpdateRepoStatus()
	gr.UpdateIndexStatus()
}

// GetLocal returns GrsRepo's directory on the local file system
func (gr *GrsRepo) GetLocal() string {
	return gr.local
}

// GetStats returns information on GrsRepo
func (gr *GrsRepo) GetStats() GrsStats {
	return gr.stats
}

// ClearError clears the error flag associated with this GrsRepo
func (gr *GrsRepo) ClearError() {
	gr.err = nil
}

// GetError returns the error flag associated with this GrsRepo
func (gr *GrsRepo) GetError() error {
	return gr.err
}
