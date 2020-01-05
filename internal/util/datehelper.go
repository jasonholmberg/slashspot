package util

import "time"

const (
	// SpotDateFormat - the default spot date formate
	SpotDateFormat = "2006-01-02"
)

// BeforeNow - is the given date before now
func BeforeNow(test string) bool {
	yN, mN, dN := time.Now().In(time.Local).Date()
	testTime, _ := time.ParseInLocation(SpotDateFormat, test, time.Local)
	yT, mT, dT := testTime.Date()
	if yN == yT && mN == mT && dN == dT {
		return false
	}
	if yN > yT || mN > mT || dN > dT {
		return true
	}
	return false
}

// AfterNow - is the given date after now
func AfterNow(test string) bool {
	yN, mN, dN := time.Now().In(time.Local).Date()
	testTime, _ := time.ParseInLocation(SpotDateFormat, test, time.Local)
	yT, mT, dT := testTime.Date()
	if yN == yT && mN == mT && dN == dT {
		return false
	}
	if yN < yT || mN < mT || dN < dT {
		return true
	}
	return false
}
