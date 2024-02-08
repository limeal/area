package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"fmt"
	"time"
)

// It calculates the time between two dates and returns the time in seconds
func calculateTime(req static.AreaRequest) shared.AreaResponse {

	// Calculate the time between two dates
	// Return the time in seconds

	from, err := time.Parse("2006-01-02T15:04:05", (*req.Store)["req:calculate:time:from"].(string))
	if err != nil {
		return shared.AreaResponse{Error: err}
	}
	to, errr := time.Parse("2006-01-02T15:04:05", (*req.Store)["req:calculate:time:to"].(string))
	if errr != nil {
		return shared.AreaResponse{Error: err}
	}

	diff := to.Sub(from)

	result := 0.0
	unit := "second"
	if (*req.Store)["req:diff:unit"] != nil {
		unit = (*req.Store)["req:diff:unit"].(string)
	}
	switch unit {
	case "minute":
		result = diff.Minutes()
	case "hour":
		result = diff.Hours()
	case "day":
		result = diff.Hours() / 24
	case "week":
		result = diff.Hours() / 24 / 7
	case "month":
		result = diff.Hours() / 24 / 30
	case "year":
		result = diff.Hours() / 24 / 365
	default:
		result = diff.Seconds()
	}

	req.Logger.WriteInfo(fmt.Sprintf("Calculated time between %v and %v: %v", from.String(), to.String(), result), false)
	return shared.AreaResponse{
		Error: nil,
	}
}

// `calculateTime` is a function that takes a `map[string]interface{}` and returns a
// `map[string]interface{}` and an `error`
func DescriptorForTimeReactionCalculateTimeDifference() static.ServiceArea {
	return static.ServiceArea{
		Name:        "time_reaction_calculate_time",
		Description: "Calculate the time between two dates",
		RequestStore: map[string]static.StoreElement{
			"req:calculate:time:from": {
				Description: "The date to calculate the time from formatted as YYYY-MM-DDTHH:MM:SS",
				Type:        "date",
				Required:    true,
			},
			"req:calculate:time:to": {
				Priority:    1,
				Description: "The date to calculate the time to formatted as YYYY-MM-DDTHH:MM:SS",
				Type:        "date",
				Required:    true,
			},
			"req:time:zone": {
				Priority:    3,
				Description: "The time zone to use (default: UTC)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"UTC-12",
					"UTC-11",
					"UTC-10",
					"UTC-9",
					"UTC-8",
					"UTC-7",
					"UTC-6",
					"UTC-5",
					"UTC-4",
					"UTC-3",
					"UTC-2",
					"UTC-1",
					"UTC",
					"UTC+1",
					"UTC+2",
					"UTC+3",
					"UTC+4",
					"UTC+5",
					"UTC+6",
					"UTC+7",
					"UTC+8",
					"UTC+9",
					"UTC+10",
					"UTC+11",
					"UTC+12",
					"UTC+13",
					"UTC+14",
				},
			},
			"req:diff:unit": {
				Priority:    2,
				Description: "The unit to return the time in (default: seconds)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"second",
					"minute",
					"hour",
					"day",
					"week",
					"month",
					"year",
				},
			},
		},
		Method: calculateTime,
	}
}
