package script

import (
	"fmt"
	"time"
)

func AutoPushGenCommitMsg(clock Clock) string {
	return fmt.Sprintf("grs-autocommit:%v", clock.Now().Format(time.RFC3339))
}

// Clock interface allows one to mock functions of the time.Time type
type Clock interface {
	Now() time.Time
}
type StdClock struct{}

func (s *StdClock) Now() time.Time {
	return time.Now()
}
func NewStdClock() *StdClock {
	return &StdClock{}
}

type MockClock struct {
	NowRetval time.Time
}

func (s *MockClock) Now() time.Time {
	return s.NowRetval
}
