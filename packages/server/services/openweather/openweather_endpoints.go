package openweather

import (
	"area-server/classes/static"
	"area-server/utils"
	"os"
)

// It returns a map of endpoint names to endpoint definitions
func OpenWeatherEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Validators
		"GetStationEndpoint": {
			BaseURL:        "https://api.openweathermap.org/data/3.0/stations/${station_id}",
			Params:         GetStationEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"GetAllStationsEndpoint": {
			BaseURL:        "https://api.openweathermap.org/data/3.0/stations",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetCoordsByLocationEndpoint": {
			BaseURL:        "http://api.openweathermap.org/geo/1.0/direct",
			Params:         GetCoordsByLocationEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetWeatherAtLocationEndpoint": {
			BaseURL:        "https://api.openweathermap.org/data/2.5/weather",
			Params:         GetWeatherAtLocationEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"CreateStationEndpoint": {
			BaseURL:        "https://api.openweathermap.org/data/3.0/stations",
			Params:         CreateStationEndpointParams,
			ExpectedStatus: []int{201},
		},
		"DeleteStationEndpoint": {
			BaseURL:        "https://api.openweathermap.org/data/3.0/stations/${station_id}",
			Params:         DeleteStationEndpointParams,
			ExpectedStatus: []int{204},
		},
	}
}

// ------------------- Endpoints --------------------------------

// It returns a pointer to a RequestParams struct with the Method, Headers, and QueryParams fields set
// to the values specified
func BasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		QueryParams: map[string]string{
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
		},
	}
}

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func GetStationEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		QueryParams: map[string]string{
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
		},
		UrlParams: map[string]string{
			"station_id": params[0].(string),
		},
	}
}

// It returns a pointer to a `RequestParams` struct that contains the necessary information to make a
// GET request to the OpenWeatherMap API
func GetCoordsByLocationEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		QueryParams: map[string]string{
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
			"q":     params[0].(string),
		},
	}
}

// It takes a slice of interfaces, and returns a pointer to a RequestParams struct
func GetWeatherAtLocationEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		QueryParams: map[string]string{
			"lat":   params[0].([]string)[0],
			"lon":   params[0].([]string)[1],
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
		},
	}
}

// It creates a request params object that will be used to make a POST request to the OpenWeatherMap
// API
func CreateStationEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Accept":      "application/json",
			"Cotent-Type": "application/json",
		},
		QueryParams: map[string]string{
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
		},
		Body: params[0].(string),
	}
}

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func DeleteStationEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "DELETE",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		QueryParams: map[string]string{
			"appId": os.Getenv("OPEN_WEATHER_API_KEY"),
		},
		UrlParams: map[string]string{
			"station_id": params[0].(string),
		},
	}
}
