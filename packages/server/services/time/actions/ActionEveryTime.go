package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/time/common"
	"area-server/utils"
	"strconv"
	"time"
)

// It returns true if the current time matches the parameters passed in
func everyTime(req static.AreaRequest) shared.AreaResponse {
	if (*req.Store)["req:year"] != nil {
		var err error
		(*req.Store)["ctx:year"], err = strconv.Atoi((*req.Store)["req:year"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:month:of:year"] != nil {
		var err error
		(*req.Store)["ctx:month"], err = utils.GetMonth((*req.Store)["req:month:of:year"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:day:of:week"] != nil {
		var err error
		(*req.Store)["ctx:day"], err = utils.GetDay((*req.Store)["req:day:of:week"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:hour:of:day"] != nil {
		var err error
		(*req.Store)["ctx:hour"], err = strconv.Atoi((*req.Store)["req:hour:of:day"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:minute:of:hour"] != nil {
		var err error
		(*req.Store)["ctx:minute"], err = strconv.Atoi((*req.Store)["req:minute:of:hour"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:time:zone"] != nil {
		var err error
		(*req.Store)["ctx:time:zone"], err = common.TransformTimeZoneIntoFixedZone((*req.Store)["req:time:zone"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	ctime := time.Now()
	if (*req.Store)["ctx:time:zone"] != nil {
		ctime = time.Now().In((*req.Store)["ctx:time:zone"].(*time.Location))
	}
	year := ctime.Year()
	month := ctime.Month()
	day := ctime.Weekday()
	hour := ctime.Hour()
	minute := ctime.Minute()
	hasParams := (*req.Store)["ctx:year"] != nil || (*req.Store)["ctx:month"] != nil || (*req.Store)["ctx:day"] != nil || (*req.Store)["ctx:hour"] != nil || (*req.Store)["ctx:minute"] != nil

	if hasParams {

		if (*req.Store)["ctx:year"] != nil && (*req.Store)["ctx:year"].(int) != year {
			return shared.AreaResponse{Success: false}
		}
		if (*req.Store)["ctx:month"] != nil && (*req.Store)["ctx:month"].(time.Month) != month {
			return shared.AreaResponse{Success: false}
		}
		if (*req.Store)["ctx:day"] != nil && (*req.Store)["ctx:day"].(time.Weekday) != day {
			return shared.AreaResponse{Success: false}
		}
		if (*req.Store)["ctx:hour"] != nil && (*req.Store)["ctx:hour"].(int) != hour {
			return shared.AreaResponse{Success: false}
		}
		if (*req.Store)["ctx:minute"] != nil && (*req.Store)["ctx:minute"].(int) != minute {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:year:passed"] != nil && (*req.Store)["ctx:year:passed"].(int) == year {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:month:passed"] != nil && (*req.Store)["ctx:month:passed"].(time.Month) == month {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:day:passed"] != nil && (*req.Store)["ctx:day:passed"].(time.Weekday) == day {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:hour:passed"] != nil && (*req.Store)["ctx:hour:passed"].(int) == hour {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:minute:passed"] != nil && (*req.Store)["ctx:minute:passed"].(int) == minute {
			return shared.AreaResponse{Success: false}
		}

		(*req.Store)["ctx:year:passed"] = year
		(*req.Store)["ctx:month:passed"] = month
		(*req.Store)["ctx:day:passed"] = day
		(*req.Store)["ctx:hour:passed"] = hour
		(*req.Store)["ctx:minute:passed"] = minute
	} else {
		if hour != 0 || minute != 0 {
			return shared.AreaResponse{Success: false}
		}

		if (*req.Store)["ctx:past:day"] != nil && (*req.Store)["ctx:past:day"].(time.Weekday) == day {
			return shared.AreaResponse{Success: false}
		}

		(*req.Store)["ctx:past:day"] = day
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"time:current:date":   time.Now().Format("2006-01-02"),
			"time:current:time":   time.Now().Format("15:04:05"),
			"time:current:zone":   time.Now().Format("-07:00"),
			"time:current:year":   time.Now().Year(),
			"time:current:month":  time.Now().Month().String(),
			"time:current:day":    time.Now().Weekday().String(),
			"time:current:hour":   time.Now().Hour(),
			"time:current:minute": time.Now().Minute(),
		},
	}
}

// `DescriptorForTimeActionEveryTime` returns a `static.ServiceArea` with the name `every_time`, a
// description, a request store, a method, and a list of components
func DescriptorForTimeActionEveryTime() static.ServiceArea {
	return static.ServiceArea{
		Name:        "every_time",
		Description: "Triggered every matching time, if no time is specified, it will be triggered every day at midnight",
		RequestStore: map[string]static.StoreElement{
			"req:year": {
				Description: "The year to trigger the action (default: every year)", // Year must be >= current year
				Type:        "string",
				Required:    false,
			},
			"req:month:of:year": {
				Priority:    1,
				Description: "The month of the year to trigger the action (default: every month)", // Month must be >= current month if year is current year
				Type:        "select",
				Required:    false,
				Values: []string{
					"January",
					"February",
					"March",
					"April",
					"May",
					"June",
					"July",
					"August",
					"September",
					"October",
					"November",
					"December",
				},
			},
			"req:day:of:week": {
				Priority:    2,
				Description: "The day of the week to trigger the action (default: every day)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"Monday",
					"Tuesday",
					"Wednesday",
					"Thursday",
					"Friday",
					"Saturday",
					"Sunday",
				},
			},
			"req:hour:of:day": {
				Priority:    3,
				Description: "The hour of the day to trigger the action (default: every hour)",
				Type:        "string",
				Required:    false,
			},
			"req:minute:of:hour": {
				Priority:    4,
				Description: "The minute of the hour to trigger the action (default: every minute)",
				Type:        "string",
				Required:    false,
			},
			"req:time:zone": {
				Priority:    5,
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
		},
		Method: everyTime,
		Components: []string{
			"time:current:date",
			"time:current:time",
			"time:current:zone",
			"time:current:year",
			"time:current:month",
			"time:current:day",
			"time:current:hour",
			"time:current:minute",
		},
	}
}
