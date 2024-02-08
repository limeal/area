package utils

import (
	"errors"
	"time"
)

// It takes a string and returns a time.Month and an error
func GetMonth(month string) (time.Month, error) {
	switch month {
	case "January":
		return time.January, nil
	case "February":
		return time.February, nil
	case "March":
		return time.March, nil
	case "April":
		return time.April, nil
	case "May":
		return time.May, nil
	case "June":
		return time.June, nil
	case "July":
		return time.July, nil
	case "August":
		return time.August, nil
	case "September":
		return time.September, nil
	case "October":
		return time.October, nil
	case "November":
		return time.November, nil
	case "December":
		return time.December, nil
	}
	return 0, errors.New("Invalid month")
}

// "Given a string, return a time.Weekday and an error."
//
// The first thing to notice is that the function returns two values. The first is a time.Weekday,
// which is an integer type. The second is an error, which is a built-in interface type
func GetDay(day string) (time.Weekday, error) {
	switch day {
	case "Monday":
		return time.Monday, nil
	case "Tuesday":
		return time.Tuesday, nil
	case "Wednesday":
		return time.Wednesday, nil
	case "Thursday":
		return time.Thursday, nil
	case "Friday":
		return time.Friday, nil
	case "Saturday":
		return time.Saturday, nil
	case "Sunday":
		return time.Sunday, nil
	}
	return 0, errors.New("Invalid day")
}
