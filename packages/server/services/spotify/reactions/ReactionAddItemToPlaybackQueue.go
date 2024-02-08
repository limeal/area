package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
)

// It adds an item to the user's playback queue
func addItemToQueue(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)
	uri := utils.GenerateFinalComponent((*req.Store)["req:uri"].(string), req.ExternalData, []string{
		"spotify:track:uri",
		"spotify:episode:uri",
	})

	query["uri"] = uri

	if (*req.Store)["req:device:id"] != nil {
		query["device_id"] = (*req.Store)["req:device:id"].(string)
	}

	_, httpResp, err := req.Service.Endpoints["AddItemToPlaybackQueue"].CallEncode([]interface{}{req.Authorization, query})

	if httpResp != nil && httpResp.StatusCode == 403 {
		req.Logger.WriteInfo("(Spotify) No item added to queue, the user is not premium", true)
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{Error: nil}
}

// `DescriptorForSpotifyReactionAddItemToPlaybackQueue` returns a `static.ServiceArea` that describes
// the `add_item_to_playback_queue` service area
func DescriptorForSpotifyReactionAddItemToPlaybackQueue() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_item_to_playback_queue",
		Description: "Add an item to the end of the user’s current playback queue. (Premium account required)",
		RequestStore: map[string]static.StoreElement{
			"req:uri": {
				Type:              "string",
				Description:       "The uri of the item to add to the queue. Must be a track or episode uri.",
				Required:          true,
				AllowedComponents: []string{"spotify:track:uri", "spotify:episode:uri"},
			},
			"req:device:id": {
				Priority:    1,
				Description: "The id of the device this command is targeting. If not supplied, the user’s currently active device is the target.",
				Required:    false,
				Type:        "select_uri",
				Values:      []string{"/devices"},
			},
		},
		Method: addItemToQueue,
	}
}
