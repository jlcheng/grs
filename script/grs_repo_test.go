package script

import (
	"jcheng/grs/shexec"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestUpdateDirstat(t *testing.T) {
	const testID = "TestUpdateDirstat"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper()
	gh.NewRepoPair(tmpdir)

	commandRunner := &shexec.ExecRunner{}
	gr := NewGrsRepo(WithCommandRunnerGrsRepo(commandRunner))
	gr.UpdateDirstat()
	if gr.stats.Dir != GRSDIR_INVALID {
		t.Fatal("unexpected Dirstat:", gr.stats.Dir)
	}
	if gr.GetError() == nil {
		t.Fatal("missing expected error")
	}

	gr = NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner))
	gr.UpdateDirstat()
	if gr.stats.Dir != GRSDIR_VALID {
		t.Fatal("unexpected Dirstat:", gr.stats.Dir)
	}
}

func TestGrsStats(t *testing.T) {
	const testID = "TestUpdateGrsStats"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper()
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}
	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner))

	gr.Update()
	expected := NewGrsStats(
		WithBranchstat(BRANCH_UPTODATE),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	gh.Touch("change-1.txt")
	gr.Update()
	expected = NewGrsStats(
		WithBranchstat(BRANCH_UPTODATE),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_MODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	gh.Add(".")
	gh.Commit("A")
	gr.Update()
	expected = NewGrsStats(
		WithBranchstat(BRANCH_AHEAD),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	gh.RunGit("push")
	gr.Update()
	expected = NewGrsStats(
		WithBranchstat(BRANCH_UPTODATE),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	gh.RunGit("reset", "--hard", "HEAD~1")
	gr.Update()
	expected = NewGrsStats(
		WithBranchstat(BRANCH_BEHIND),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	gh.Touch("change-2.txt")
	gh.Add(".")
	gh.Commit("B")
	gr.Update()
	expected = NewGrsStats(
		WithBranchstat(BRANCH_DIVERGED),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}
}

func TestFetch(t *testing.T) {
	const testID = "TestFetch"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper()
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	cloneA := path.Join(tmpdir, "dest")
	cloneB := path.Join(tmpdir, "dest-b")
	gh.Chdir(tmpdir)
	gh.RunGit("clone", "source", "dest-b")
	gh.Chdir(cloneB)
	gh.TouchAndCommit("from-b.txt", "B")
	gh.RunGit("push")
	gh.Chdir(cloneA)
	gr := NewGrsRepo(WithLocalGrsRepo(cloneA), WithCommandRunnerGrsRepo(commandRunner))
	gr.Update()
	gr.Fetch()
	gr.Update()

	expected := NewGrsStats(
		WithBranchstat(BRANCH_BEHIND),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}
}

func TestUpdateCommitTime(t *testing.T) {
	const testID = "TestUpdateCommitTime"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper()
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner))
	gr.Update()
	if !strings.HasSuffix(gr.GetStats().CommitTime, " ago") {
		t.Fatal("unexpected CommitTime:", gr.GetStats().CommitTime)
	}
}

func TestAutoPush(t *testing.T) {
	const testID = "TestAutoPush"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper()
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner), WithPushAllowed(true))
	gh.TouchAndCommit("new.txt", "A")
	gr.Update()
	gr.AutoPush()

	expected := NewGrsStats(
		WithBranchstat(BRANCH_UPTODATE),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}
}

func TestAutoPushUntracked(t *testing.T) {
	const testID = "TestAutoPushUntracked"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner), WithPushAllowed(true))
	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("push")
	gh.RunGit("reset", "--hard", "A")
	gh.Touch("C.txt")

	gr.Update()
	gr.AutoPush()

	expected := NewGrsStats(
		WithBranchstat(BRANCH_DIVERGED),
		WithDirstat(GRSDIR_VALID),
		WithIndexstat(INDEX_UNMODIFIED),
	)
	if noCommitTime(gr.GetStats()) != expected {
		t.Fatal("unexpected stats:", gr.GetStats())
	}

	logGraph, _ := gh.LogGraph()
	if _, ok := logGraph["A"]; !ok {
		t.Fatal("missing commit A")
	}
	var commit = false
	for k, v := range logGraph {
		if strings.HasPrefix(k, "grs-autocommit") && v[0] == "A" {
			commit = true
		}
	}
	if !commit {
		t.Fatal("missing autocommit")
	}
}

func TestAutoRebase(t *testing.T) {
	const testID = "TestAutoRebase"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner))

	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("tag", "B")
	gh.TouchAndCommit("C.txt", "C")
	gh.RunGit("checkout", "-b", "source_2", "B")
	gh.TouchAndCommit("D.txt", "D")
	gh.TouchAndCommit("E.txt", "E")
	gh.RunGit("checkout", "master")
	gh.RunGit("merge", "source_2", "-m", "F")
	gh.RunGit("push")
	gh.RunGit("reset", "--hard", "A")
	gh.TouchAndCommit("G.txt", "G")
	gh.TouchAndCommit("H.txt", "H")

	gr.Update()
	gr.AutoRebase()
	gr.Update()

	got := LogGraph(map[string][]string{
		"A":    {"init"},
		"B":    {"A"},
		"C":    {"B"},
		"D":    {"B"},
		"E":    {"D"},
		"F":    {"C", "E"},
		"G":    {"F"},
		"H":    {"G"},
		"init": {},
	})
	expected, _ := gh.LogGraph()
	if !reflect.DeepEqual(got, expected) {
		t.Fatal("unexpected commit graph", got)
	}
}

func TestAutoFFMerge(t *testing.T) {
	const testID = "TestAutoRebase"
	tmpdir, cleanup := MkTmpDir(t, testID)
	defer cleanup()
	gh := NewGitTestHelper(WithDebug(false), WithWd(tmpdir))
	gh.NewRepoPair(tmpdir)
	commandRunner := &shexec.ExecRunner{}

	gr := NewGrsRepo(WithLocalGrsRepo(gh.Getwd()), WithCommandRunnerGrsRepo(commandRunner))

	gh.TouchAndCommit("A.txt", "A")
	gh.RunGit("tag", "A")
	gh.TouchAndCommit("B.txt", "B")
	gh.RunGit("push")
	gh.RunGit("reset", "--hard", "A")

	gr.Update()
	gr.AutoFFMerge()
	gr.Update()

	got := LogGraph(map[string][]string{
		"A":    {"init"},
		"B":    {"A"},
		"init": {},
	})
	expected, _ := gh.LogGraph()
	if !reflect.DeepEqual(got, expected) {
		t.Fatal("unexpected commit graph", got)
	}
}

func noCommitTime(stats GrsStats) GrsStats {
	copy := stats
	copy.CommitTime = ""
	return copy
}
