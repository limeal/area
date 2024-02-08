package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
	"fmt"
)

// GetCoordsResponse is a struct that contains a slice of structs that each contain a string, a
// float64, and a float64.
// @property {[]struct {
// 		Name      string  `json:"name"`
// 		Latitude  float64 `json:"latitude"`
// 		Longitude float64 `json:"longitude"`
// 	}} Cities - An array of cities.
type GetCoordsResponse struct {
	Cities []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
}

// Weather is a struct that contains two fields, Main and Description, both of which are strings.
// @property {string} Main - This is the main weather condition. For example, "Rain" or "Clouds".
// @property {string} Description - The description of the weather.
type Weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

// GetWeatherResponse is a struct that contains a Coordinate, a slice of Weather, a timestamp, and an
// ID.
// @property Coords - The latitude and longitude of the location
// @property {[]Weather} Weather - This is an array of weather objects.
// @property {int} DT - Time of data calculation, unix, UTC
// @property {int} ID - City ID
type GetWeatherResponse struct {
	Coords struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"coord"`
	Weather []Weather `json:"weather"`
	DT      int       `json:"dt"`
	ID      int       `json:"id"`
}

// It checks if the weather is raining at a given location
func isRainingAtLocation(req static.AreaRequest) shared.AreaResponse {

	query := ""
	if (*req.Store)["ctx:location"] == nil {
		if (*req.Store)["req:country:code"] == nil {
			query = fmt.Sprintf("%s,%s", (*req.Store)["req:city:name"], (*req.Store)["req:country:code"])
		} else {
			query = (*req.Store)["req:city:name"].(string)
		}

		encode, _, err := req.Service.Endpoints["GetCoordsByLocationEndpoint"].CallEncode([]interface{}{
			query,
		})

		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		cities := GetCoordsResponse{}
		errr := json.Unmarshal(encode, &cities.Cities)
		if errr != nil {
			return shared.AreaResponse{Error: err}
		}

		(*req.Store)["ctx:location"] = []string{
			fmt.Sprintf("%v", cities.Cities[0].Latitude),
			fmt.Sprintf("%v", cities.Cities[0].Longitude),
		}
	}

	encode, _, err := req.Service.Endpoints["GetWeatherAtLocationEndpoint"].CallEncode([]interface{}{
		(*req.Store)["ctx:location"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var body GetWeatherResponse
	if err := json.Unmarshal(encode, &body); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbWeather := len(body.Weather)
	if (*req.Store)["ctx:total_weather"] == nil {
		(*req.Store)["ctx:total_weather"] = nbWeather
	}

	if (*req.Store)["ctx:total_weather"].(int) == 0 && nbWeather == 0 {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:total_weather"].(int) > 0 && (*req.Store)["ctx:main:weather"] == nil {
		(*req.Store)["ctx:main:weather"] = body.Weather[0].Main
	}

	if (*req.Store)["ctx:total_weather"].(int) > 0 && (*req.Store)["ctx:main:weather"].(string) == body.Weather[0].Main {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"openw:weather:id":          body.ID,
			"openw:weather:name":        body.Weather[0].Main,
			"openw:weather:description": body.Weather[0].Description,
			"openw:coords:latitude":     body.Coords.Latitude,
			"openw:coords:longitude":    body.Coords.Longitude,
			"openw:city:name":           (*req.Store)["req:city:name"],
			"openw:country:name":        (*req.Store)["req:country:code"],
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForOpenWeatherActionWeatherChangeAtLocation() static.ServiceArea {
	return static.ServiceArea{
		Name:        "weather_change_at_location",
		Description: "Triggered if the weather change at a specific location",
		RequestStore: map[string]static.StoreElement{
			"req:country:code": {
				Description: "You can specify the country code to handle a more specific request",
				Type:        "string",
				Required:    false,
			},
			"req:city:name": {
				Priority:    1,
				Description: "The name of the city to specify if it's raining",
				Type:        "string",
				Required:    true,
			},
			"req:wheather:type": {
				Priority:    2,
				Description: "The type of weather you want to check",
				Type:        "select",
				Required:    false,
				Values: []string{
					"Thunderstorm",
					"Drizzle",
					"Rain",
					"Snow",
					"Mist",
					"Smoke",
					"Haze",
					"Dust",
					"Fog",
					"Sand",
					"Ash",
					"Squall",
					"Tornado",
					"Clear",
					"Clouds",
				},
			},
		},
		Method: isRainingAtLocation,
		Components: []string{
			"openw:weather:id",
			"openw:weather:name",
			"openw:weather:description",
			"openw:coords:latitude",
			"openw:coords:longitude",
			"openw:city:name",
			"openw:country:name",
		},
	}
}
