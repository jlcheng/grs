package script

import (
	"fmt"
	"time"
)

// AutoPushGenCommitMsg creates a commit message from a Clock instance
func AutoPushGenCommitMsg(clock Clock) string {
	return fmt.Sprintf("grs-autocommit:%v", clock.Now().Format(time.RFC3339))
}

// Clock interface allows one to mock functions of the time.Time type
type Clock interface {
	Now() time.Time
}

// StdClock implements Clock using time.Time
type StdClock struct{}

// Now returns the current time
func (s *StdClock) Now() time.Time {
	return time.Now()
}

// NewStdClock returns a pointer to a new StdClock instance
func NewStdClock() *StdClock {
	return &StdClock{}
}

// MockClock is a wide-open implementation of Clock; Anyone can modify its "Now" time.
type MockClock struct {
	NowRetval time.Time
}

// Now returns the value of the NowRetVal field
func (s *MockClock) Now() time.Time {
	return s.NowRetval
}
