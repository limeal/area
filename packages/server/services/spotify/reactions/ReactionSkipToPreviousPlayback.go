package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
)

// It skips to the previous song in the current playlist
func skipToPreviousPlayback(req static.AreaRequest) shared.AreaResponse {

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

	_, httpResp, err := req.Service.Endpoints["SkipToPreviousPlaybackEndpoint"].CallEncode([]interface{}{req.Authorization, device})

	if httpResp != nil && httpResp.StatusCode == 403 {
		req.Logger.WriteInfo("(Spotify) Playback can't be skipped to previous, the user is not premium", true)
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{Error: nil}
}

// `DescriptorForSpotifyReactionSkipToPreviousPlayback` returns a `static.ServiceArea` with the name
// `skip_to_previous_playback`, a description of `Skip to previous playback song (Premium account
// required)`, a request store with a single element `req:device:id` that is not required, is of type
// `select_uri` and has a value of `/devices`, and a method of `skipToPreviousPlayback`
func DescriptorForSpotifyReactionSkipToPreviousPlayback() static.ServiceArea {
	return static.ServiceArea{
		Name:        "skip_to_previous_playback",
		Description: "Skip to previous playback song (Premium account required)",
		RequestStore: map[string]static.StoreElement{
			"req:device:id": {
				Description: "The id of the device this command is targeting. If not supplied, the user’s currently active device is the target.",
				Required:    false,
				Type:        "select_uri",
				Values:      []string{"/devices"},
			},
		},
		Method: skipToPreviousPlayback,
	}
}
