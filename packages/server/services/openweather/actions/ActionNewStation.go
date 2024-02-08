package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/openweather/common"
	"area-server/utils"
	"encoding/json"
)

// It gets the latest station from the OpenWeatherMap API, and returns it as a response
func onNewStation(req static.AreaRequest) shared.AreaResponse {

	encode, _, err := req.Service.Endpoints["GetAllStationsEndpoint"].CallEncode([]interface{}{})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var stations []common.Station
	err = json.Unmarshal(encode, &stations)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbStations := len(stations)
	if ok, err := utils.IsLatestByDate(req.Store, nbStations, func() interface{} {
		return stations[0].ID
	}, func() string {
		return stations[0].CreatedAt
	}); err != nil || !ok {
		return shared.AreaResponse{Error: err, Success: false}
	}

	newStation := stations[0]
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"openw:station:id":          newStation.ID,
			"openw:station:name":        newStation.Name,
			"openw:station:external:id": newStation.ExternalID,
			"openw:station:latitude":    newStation.Latitude,
			"openw:station:longitude":   newStation.Longitude,
			"openw:station:altitude":    newStation.Altitude,
			"openw:station:rank":        newStation.Rank,
		},
	}
}

// It returns a `static.ServiceArea` struct that describes the action `onNewStation` and the components
// it uses
func DescriptorForOpenWeatherActionAnyNewStation() static.ServiceArea {
	return static.ServiceArea{
		Name:        "any_new_station",
		Description: "Check if a new station has been added",
		WIP:         true,
		Method:      onNewStation,
		Components: []string{
			"openw:station:id",
			"openw:station:name",
			"openw:station:external:id",
			"openw:station:latitude",
			"openw:station:longitude",
			"openw:station:altitude",
			"openw:station:rank",
		},
	}
}
