package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// `DevicesResponse` is a struct with a field `Devices` which is a slice of structs with fields `ID`,
// `IsActive`, `IsPrivateSession`, `IsRestricted`, `Name`, `Type`, and `VolumePercent`.
// @property {[]struct {
// 		ID               string `json:"id"`
// 		IsActive         bool   `json:"is_active"`
// 		IsPrivateSession bool   `json:"is_private_session"`
// 		IsRestricted     bool   `json:"is_restricted"`
// 		Name             string `json:"name"`
// 		Type             string `json:"type"`
// 		VolumePercent    int    `json:"volume_percent"`
// 	}} Devices - An array of devices.
type DevicesResponse struct {
	Devices []struct {
		ID               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
	} `json:"devices"`
}

// It resumes the playback of the current user
func resumePlayback(req static.AreaRequest) shared.AreaResponse {

	body := map[string]interface{}{}
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

	if (*req.Store)["req:context:uri"] != nil {
		body["context_uri"] = utils.GenerateFinalComponent((*req.Store)["req:context:uri"].(string), req.ExternalData, []string{
			"spotify:album:uri",
			"spotify:artist:uri",
		})
	}

	if (*req.Store)["req:element:uri"] != nil {
		body["uris"] = []string{utils.GenerateFinalComponent((*req.Store)["req:element:uri"].(string), req.ExternalData, []string{
			"spotify:track:uri",
		})}
	}

	ebody, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, httpResp, errr := req.Service.Endpoints["ResumePlaybackEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		device,
		string(ebody),
	})

	if httpResp != nil && httpResp.StatusCode == 403 {
		req.Logger.WriteInfo("(Spotify) Playback can't be resume, the user is not premium", true)
	}

	if errr != nil {
		return shared.AreaResponse{Error: err}
	}
	return shared.AreaResponse{Error: nil}
}

// `DescriptorForSpotifyReactionResumePlayback` returns a `static.ServiceArea` that describes the
// `resumePlayback` function
func DescriptorForSpotifyReactionResumePlayback() static.ServiceArea {
	return static.ServiceArea{
		Name:        "start_playback",
		Description: "Start playback (Premium account required)",
		RequestStore: map[string]static.StoreElement{
			"req:device:id": {
				Priority:    1,
				Description: "The id of the device this command is targeting. If not supplied, the user’s currently active device is the target.",
				Required:    false,
				Type:        "select_uri",
				Values:      []string{"/devices"},
			},
			"req:context:uri": {
				Priority:          2,
				Type:              "string",
				Description:       "The context to play. Valid contexts are albums, artists & playlists.",
				Required:          false,
				AllowedComponents: []string{"spotify:album:uri", "spotify:artist:uri"},
			},
			"req:element:uri": {
				Priority:          1,
				Type:              "string",
				Description:       "The uri of the element to play. Can be a track or an episode uri.",
				Required:          false,
				AllowedComponents: []string{"spotify:track:uri"},
			},
		},
		Method: resumePlayback,
	}
}
