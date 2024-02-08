package common

import (
	"strconv"
	"time"
)

// zone is formatted as UTC+00
// It converts a timezone string like "+0100" into a time.Location object
func TransformTimeZoneIntoFixedZone(formatZone string) (*time.Location, error) {
	fsign := formatZone[3]
	hours := formatZone[4:]
	sign := 1
	if fsign == '-' {
		sign = -1
	}
	// Convert to int
	hoursInt, errh := strconv.Atoi(hours)
	if errh != nil {
		return nil, errh
	}
	// Convert to seconds
	hoursInt *= 60 * 60
	return time.FixedZone(formatZone, hoursInt*sign), nil
}
