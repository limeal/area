package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"strconv"
	"time"
)

// It waits for a specified amount of time, and then returns the current date and time, the target date
// and time, and the duration and unit of time that was waited
func waitTime(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["ctx:time:duration"] == nil {
		var err error
		(*req.Store)["ctx:time:duration"], err = strconv.ParseFloat((*req.Store)["req:time:duration"].(string), 64)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["ctx:end:time"] == nil {
		switch (*req.Store)["req:time:unit"] {
		case "minute":
			(*req.Store)["ctx:time:unit"] = "minute"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * time.Minute)
		case "hour":
			(*req.Store)["ctx:time:unit"] = "hour"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * time.Hour)
		case "day":
			(*req.Store)["ctx:time:unit"] = "day"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * 24 * time.Hour)
		case "week":
			(*req.Store)["ctx:time:unit"] = "week"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * 7 * 24 * time.Hour)
		case "month":
			(*req.Store)["ctx:time:unit"] = "month"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * 30 * 24 * time.Hour)
		case "year":
			(*req.Store)["ctx:time:unit"] = "year"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * 365 * 24 * time.Hour)
		default:
			(*req.Store)["ctx:time:unit"] = "second"
			(*req.Store)["ctx:end:time"] = time.Now().Add(time.Duration((*req.Store)["ctx:time:duration"].(float64)) * time.Second)
		}
	}

	// Check if the time is up
	if time.Now().Before((*req.Store)["ctx:end:time"].(time.Time)) {
		return shared.AreaResponse{Success: false}
	}

	targetDate := (*req.Store)["ctx:end:time"].(time.Time)
	timeUnit := (*req.Store)["ctx:time:unit"].(string)
	(*req.Store)["ctx:end:time"] = nil
	(*req.Store)["ctx:time:unit"] = nil
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"time:wait:duration": (*req.Store)["req:time:duration"],
			"time:wait:unit":     timeUnit,
			"time:current:date":  time.Now().Format("2006-01-02-15-04-05"),
			"time:target:date":   targetDate.Format("2006-01-02-15-04-05"),
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForTimeActionWaitTime() static.ServiceArea {
	return static.ServiceArea{
		Name:        "wait_time",
		Description: "Wait for a specific amount of time",
		RequestStore: map[string]static.StoreElement{
			"req:time:duration": {
				Description: "The amount of time to wait (default: seconds)",
				Type:        "string",
				Required:    true,
			},
			"req:time:unit": {
				Priority:    1,
				Description: "The unit of time to wait",
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
		Method: waitTime,
		Components: []string{
			"time:wait:duration",
			"time:wait:unit",
			"time:current:date",
			"time:target:date",
		},
	}
}
