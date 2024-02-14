package timer

import "time"

func New() Timer {
	return &timer{}
}

type Timer interface {
	Now() time.Time
	UnixMilli(msec int64) time.Time
}

type timer struct {
}

// Now return current UTC time
func (t *timer) Now() time.Time {
	return time.Now().UTC()
}

// Now return current UTC time from milliseconds
func (t *timer) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec).UTC()
}
