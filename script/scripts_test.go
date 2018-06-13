package script

import (
	"testing"
)

func TestParseRevList(t *testing.T) {
	assertValidRemoteDiff(t, "0\t0\n", 0, 0)
	assertValidRemoteDiff(t, "1\t0\n", 1, 0)
	assertValidRemoteDiff(t, "0\t1\n", 0, 1)
	assertValidRemoteDiff(t, "1\t1\n", 1, 1)
	assertValidRemoteDiff(t, "999\t555\n", 999, 555)
	assertValidRemoteDiff(t, "555\t999\n", 555, 999)

	assertInvalidRemoteDiff(t, "\t\n")
	assertInvalidRemoteDiff(t, "555999\n")
	assertInvalidRemoteDiff(t, "555 999\n")
	assertInvalidRemoteDiff(t, "555        999\n")
	assertInvalidRemoteDiff(t, "555t999\n")
}

func assertValidRemoteDiff(t *testing.T, str string, eRemote int, eLocal int) {
	var d RemoteDiff
	var err error
	d, err = parseRevList([]byte(str))
	if err != nil {
		t.Errorf("cannot parse %v\n", str)
	}
	if d.remote != eRemote {
		t.Errorf("from %v, expected remote=%v, got remote=%v", str, eRemote, d.remote)
	}
	if d.local != eLocal {
		t.Errorf("from %v, expected local=%v, got local=%v", str, eLocal, d.local)
	}
}

func assertInvalidRemoteDiff(t *testing.T, str string) {
	var err error
	_, err = parseRevList([]byte(str))
	if err == nil {
		t.Error("expected error parsing", err)
	}
}
