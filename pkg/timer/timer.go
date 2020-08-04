package timer

import "time"

type SecondsTimer struct {
    end time.Time
}

func NewSecondsTimer(t time.Duration) *SecondsTimer {
    return &SecondsTimer{
        end: time.Now().Add(t),
    }
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
    return s.end.Sub(time.Now())
}
