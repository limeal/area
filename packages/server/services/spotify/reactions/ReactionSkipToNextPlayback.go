package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
)

// It skips to the next song in the current playlist
func skipToNextPlayback(req static.AreaRequest) shared.AreaResponse {

	device := ""
	if (*req.Store)["req:device:id"] != nil {
		device = (*req.Store)["req:device:id"].(string)
	} else {
		// Take the first available device
		encode, _, err := req.Service.Endpoints["GetUserAvailableDevicesEndpoint"].CallEncode([]interface{}{req.Authorization})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		devices := DevicesResponse{}
		err = json.Unmarshal(encode, &devices)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if len(devices.Devices) == 0 {
			req.Logger.WriteInfo("No device available", false)
			return shared.AreaResponse{Error: nil}
		}

		device = devices.Devices[0].ID
	}

	_, httpResp, err := req.Service.Endpoints["SkipToNextPlaybackEndpoint"].CallEncode([]interface{}{req.Authorization, device})

	if httpResp != nil && httpResp.StatusCode == 403 {
		req.Logger.WriteInfo("(Spotify) Playback can't be skipped to next, the user is not premium", true)
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{Error: nil}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForSpotifyReactionSkipToNextPlayback() static.ServiceArea {
	return static.ServiceArea{
		Name:        "skip_to_next_playback",
		Description: "Skip to next playback song (Premium account required)",
		RequestStore: map[string]static.StoreElement{
			"req:device:id": {
				Description: "The id of the device this command is targeting. If not supplied, the user’s currently active device is the target.",
				Required:    false,
				Type:        "select_uri",
				Values:      []string{"/devices"},
			},
		},
		Method: skipToNextPlayback,
	}
}
