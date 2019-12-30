package util

import "time"

// BeforeNow - is the given date before now
func BeforeNow(test time.Time) bool {
	now := time.Now()
	return now.Year() > test.Year() || now.YearDay() > test.YearDay()

}

// AfterNow - is the given date after now
func AfterNow(test time.Time) bool {
	now := time.Now()
	return now.Year() < test.Year() || now.YearDay() < test.YearDay()
}
