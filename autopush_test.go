package grs

import (
	"strings"
	"testing"
	"time"
)

func TestAutoPushGenCommitMsg(t *testing.T) {
	nowRetval, err := time.Parse(time.RFC3339, "1234-05-06T07:08:09Z")
	if err != nil {
		t.Error(err)
	}
	clock := &MockClock{NowRetval: nowRetval}
	if got := AutoPushGenCommitMsg(clock); !strings.Contains(got, "1234-05-06T07:08:09Z") {
		t.Error("expected timestamp missing. got:", got)
	}

}
