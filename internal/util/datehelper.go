package util

import (
	"time"
)

const (
	// SpotDateFormat - the default spot date formate
	SpotDateFormat = "2006-01-02"
	oneDay         = 24 * time.Hour
)

// Need to make sure the now time is really UTC, some sytems do not
// default to UTC.  Trying to compare to local is problematic

// BeforeNow - is the given date before now
func BeforeNow(in string) bool {
	now := time.Now().In(time.UTC).Truncate(oneDay)
	test, _ := time.ParseInLocation(SpotDateFormat, in, time.UTC)
	if now.Equal(test) {
		return false
	}
	return test.Before(now)
}

// AfterNow - is the given date after now
func AfterNow(in string) bool {
	now := time.Now().In(time.UTC).Truncate(oneDay)
	test, _ := time.ParseInLocation(SpotDateFormat, in, time.UTC)
	if now.Equal(test) {
		return false
	}
	return test.After(now)
}
