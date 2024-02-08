package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
)

// It calls the DeleteStationEndpoint endpoint, passing in the station ID from the request store
func deleteStation(req static.AreaRequest) shared.AreaResponse {

	_, _, err := req.Service.Endpoints["DeleteStationEndpoint"].Call([]interface{}{
		(*req.Store)["req:station:id"],
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForOpenWeatherReactionDeleteNewStation() static.ServiceArea {
	return static.ServiceArea{
		Name:        "delete_station",
		Description: "Delete a station",
		WIP:         true,
		RequestStore: map[string]static.StoreElement{
			"req:station:id": {
				Type:        "select_uri",
				Description: "The ID of the station",
				Required:    true,
				Values:      []string{"/stations"},
			},
		},
		Method: deleteStation,
	}
}
