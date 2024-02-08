package time

import (
	"area-server/classes/static"
	"area-server/services/time/actions"
	"area-server/services/time/reactions"
)

// It returns a static.Service object that describes the service
func Descriptor() static.Service {
	return static.Service{
		Name:        "time",
		Description: "Time service",
		RateLimit:   0,
		More: &static.More{
			Avatar: true,
			Color:  "#7289DA",
		},
		Validators: TimeValidators(),
		Actions: []static.ServiceArea{
			actions.DescriptorForTimeActionWaitTime(),
			actions.DescriptorForTimeActionEveryTime(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForTimeReactionCalculateTimeDifference(),
		},
	}
}
