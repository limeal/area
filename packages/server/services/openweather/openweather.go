package openweather

import (
	"area-server/classes/static"
	"area-server/services/openweather/actions"
	"area-server/services/openweather/reactions"
)

// It returns a static.Service object that describes the service
func Descriptor() static.Service {
	return static.Service{
		Name:        "openweather",
		Description: "Open Weather is a weather service that provides weather information",
		RateLimit:   3,
		More: &static.More{
			Avatar: true,
			Color:  "#e0533c",
		},
		Validators: OpenWeatherValidators(),
		Endpoints:  OpenWeatherEndpoints(),
		Routes:     OpenWeatherRoutes(),
		Actions: []static.ServiceArea{
			actions.DescriptorForOpenWeatherActionWeatherChangeAtLocation(),
			actions.DescriptorForOpenWeatherActionAnyNewStation(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForOpenWeatherReactionCreateStation(),
			reactions.DescriptorForOpenWeatherReactionDeleteNewStation(),
		},
	}
}
