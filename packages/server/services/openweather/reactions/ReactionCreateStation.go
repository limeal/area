package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"strconv"
)

// It creates a new station in the database
func createNewStation(req static.AreaRequest) shared.AreaResponse {

	externalId := utils.GenerateFinalComponent((*req.Store)["req:external:id"].(string), req.ExternalData, []string{})
	stationName := utils.GenerateFinalComponent((*req.Store)["req:station:name"].(string), req.ExternalData, []string{})
	stationLatitude := utils.GenerateFinalComponent((*req.Store)["req:station:latitude"].(string), req.ExternalData, []string{
		".+:latitude",
	})
	stationLongitude := utils.GenerateFinalComponent((*req.Store)["req:station:longitude"].(string), req.ExternalData, []string{
		".+:longitude",
	})
	stationAltitude := utils.GenerateFinalComponent((*req.Store)["req:station:altitude"].(string), req.ExternalData, []string{
		".+:altitude",
	})

	stationLat, errC := strconv.ParseFloat(stationLatitude, 64)
	if errC != nil {
		return shared.AreaResponse{Error: errC}
	}
	stationLong, errC := strconv.ParseFloat(stationLongitude, 64)
	if errC != nil {
		return shared.AreaResponse{Error: errC}
	}
	stationAlt, errC := strconv.ParseFloat(stationAltitude, 64)
	if errC != nil {
		return shared.AreaResponse{Error: errC}
	}

	ebody, err := json.Marshal(map[string]interface{}{
		"external_id": externalId,
		"name":        stationName,
		"latitude":    stationLat,
		"longitude":   stationLong,
		"altitude":    stationAlt,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, err = req.Service.Endpoints["CreateStationEndpoint"].Call([]interface{}{
		string(ebody),
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForOpenWeatherReactionCreateStation() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_station",
		Description: "Create a new station",
		WIP:         true,
		RequestStore: map[string]static.StoreElement{
			"req:external:id": {
				Type:        "string",
				Description: "The ID of the station",
				Required:    true,
			},
			"req:station:name": {
				Priority:    1,
				Type:        "string",
				Description: "The name of the station",
				Required:    true,
			},
			"req:station:latitude": {
				Priority:    2,
				Type:        "string",
				Description: "The latitude of the station",
				Required:    true,
			},
			"req:station:longitude": {
				Priority:    3,
				Type:        "string",
				Description: "The longitude of the station",
				Required:    true,
			},
			"req:station:altitude": {
				Priority:    4,
				Type:        "string",
				Description: "The altitude of the station",
				Required:    true,
			},
		},
		Method: createNewStation,
	}
}
