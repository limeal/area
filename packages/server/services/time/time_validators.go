package time

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/common"
	commonr "area-server/services/time/common"
	"regexp"
	"strconv"
	"time"
)

// It returns a map of validators for the time service
func TimeValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:time:zone":           ZoneValidator,
		"req:time:duration":       common.IntValidator,
		"req:time:unit":           TimeUnitValidator,
		"req:date":                DateValidator,
		"req:calculate:time:from": DateValidator,
		"req:calculate:time:to":   AfterDateValidator,
		"req:year":                YearValidator,
		"req:month:of:year":       MonthValidator,
		"req:day:of:week":         DayValidator,
		"req:hour:of:day":         HourValidator,
		"req:minute:of:hour":      MinuteValidator,
	}
}

// If the value is a string, then it must match the regex pattern for a UTC timezone.
func ZoneValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Must mast +00:00 format (regex)
	if _, ok := value.(string); ok {
		// Regex
		match, err := regexp.MatchString("UTC(\\+([1-9]){1}([0-4]){0,1}|\\-([1-9]){1}([0-2]){0,1})", value.(string))
		if err != nil {
			return false
		}
		if !match {
			return false
		}
		return true
	}

	return false
}

// If the value is a string, it must be one of the following values: second, minute, hour, day, week,
// month, year.
func TimeUnitValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Must be one of these values
	if _, ok := value.(string); ok {
		if value == "second" || value == "minute" || value == "hour" || value == "day" || value == "week" || value == "month" || value == "year" {
			return true
		}
	}

	return false
}

// It checks if the value is a string, and if it is, it checks if it's a valid date
func DateValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Must be one of these values
	if _, ok := value.(string); ok {
		var ctime time.Time
		if store["req:time:zone"] != nil {
			llocation, err := commonr.TransformTimeZoneIntoFixedZone(store["req:time:zone"].(string))
			if err != nil {
				return false
			}
			ctime, err = time.ParseInLocation("2006-01-02T15:04:05", value.(string), llocation)

			if err != nil {
				return false
			}
		} else {
			var err error
			ctime, err = time.Parse("2006-01-02T15:04:05", value.(string))
			if err != nil {
				return false
			}
		}

		if ctime.IsZero() {
			return false
		}
		return true
	}

	return false
}

// `AfterDateValidator` checks if the value is a string and if it is a valid date in the future
func AfterDateValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Must be one of these values
	if _, ok := value.(string); ok {
		var ctime time.Time
		if store["req:time:zone"] != nil {
			llocation, err := commonr.TransformTimeZoneIntoFixedZone(store["req:time:zone"].(string))
			if err != nil {
				return false
			}
			ctime, err = time.ParseInLocation("2006-01-02T15:04:05", value.(string), llocation)

			if err != nil {
				return false
			}
		} else {
			var err error
			ctime, err = time.Parse("2006-01-02T15:04:05", value.(string))
			if err != nil {
				return false
			}
		}

		if ctime.IsZero() || ctime.UTC().Before(time.Now().UTC()) {
			return false
		}
		return true
	}

	return false
}

// "If the value is not nil, parse it as an integer and return true if it's greater than or equal to
// the current year."
//
// The first thing we do is check if the value is nil. If it is, we return false. If it's not, we parse
// it as an integer. If we can't parse it, we return false. If we can, we check if it's greater than or
// equal to the current year. If it is, we return true. If it's not, we return false
func YearValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Parse year
	year, err := strconv.Atoi(value.(string))
	if err != nil {
		return false
	}

	if year < time.Now().Year() {
		return false
	}

	return true
}

// "If the value is not nil, and it's a valid month, return true."
//
// The `auth` parameter is the authorization object that is being validated. The `service` parameter is
// the service that is being accessed. The `value` parameter is the value of the parameter being
// validated. The `store` parameter is a map of values that can be used to store state between
// validators
func MonthValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Parse month
	switch value.(string) {
	case "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December":
		return true
	}

	return false
}

// "If the value is not nil, and it's a valid day of the week, return true."
//
// The first thing we do is check if the value is nil. If it is, we return false
func DayValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Parse day
	switch value.(string) {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday":
		return true
	}

	return false
}

// "If the value is not nil, parse it as an integer and return true if it's between 0 and 23, otherwise
// return false."
//
// The `HourValidator` function is a validator function. It takes four arguments:
//
// * `auth`: The authorization object.
// * `service`: The service object.
// * `value`: The value to validate.
// * `store`: A map of values that can be used to store data for later use
func HourValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Parse hour
	hour, err := strconv.Atoi(value.(string))
	if err != nil {
		return false
	}

	if hour < 0 || hour > 23 {
		return false
	}

	return true
}

// "If the value is not nil, parse it as an integer and return true if it's between 0 and 59, otherwise
// return false."
//
// The `auth` parameter is the authorization being validated. The `service` parameter is the service
// being accessed. The `value` parameter is the value of the parameter being validated. The `store`
// parameter is a map that can be used to store data that can be used by other validators
func MinuteValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	// Parse minute
	minute, err := strconv.Atoi(value.(string))
	if err != nil {
		return false
	}

	if minute < 0 || minute > 59 {
		return false
	}

	return true
}
