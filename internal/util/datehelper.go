package util

import "time"

// BeforeNow - is the given date before now
func BeforeNow(test time.Time) bool {
	now := time.Now()
	return test.Truncate(24 * time.Hour).Before(now.Truncate(24 * time.Hour))
}

// AfterNow - is the given date after now
func AfterNow(test time.Time) bool {
	now := time.Now()
	return test.Truncate(24 * time.Hour).After(now.Truncate(24 * time.Hour))
}
