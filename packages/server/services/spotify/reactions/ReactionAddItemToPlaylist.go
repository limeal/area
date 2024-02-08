package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"errors"
)

// It adds an item to a playlist
func addItemToPlaylist(req static.AreaRequest) shared.AreaResponse {
	if req.ExternalData["spotify:playlist:id"] != nil && req.ExternalData["spotify:playlist:id"] == (*req.Store)["req:writable:playlist:id"] {
		return shared.AreaResponse{Error: errors.New("Cannot add to the same playlist")}
	}

	playlistID := utils.GenerateFinalComponent((*req.Store)["req:writable:playlist:id"].(string), req.ExternalData, []string{
		"spotify:playlist:id",
	})
	uri := utils.GenerateFinalComponent((*req.Store)["req:uri"].(string), req.ExternalData, []string{
		"spotify:track:id",
		"spotify:episode:id",
	})

	body := map[string]interface{}{
		"uris":     []interface{}{uri},
		"position": 0,
	}

	str, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["AddItemToPlaylistEndpoint"].CallEncode([]interface{}{req.Authorization, playlistID, string(str)})

	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	return shared.AreaResponse{Error: nil}
}

// It returns a `static.ServiceArea` object that describes the service area `add_item_to_playlist`
func DescriptorForSpotifyReactionAddItemToPlaylist() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_item_to_playlist",
		Description: "Add a track or episode to a playlist",
		RequestStore: map[string]static.StoreElement{
			"req:writable:playlist:id": {
				Type:              "select_uri",
				Description:       "ID of the playlist to add the track to",
				Required:          true,
				AllowedComponents: []string{"spotify:playlist:id"},
				Values:            []string{"/playlists"},
			},
			"req:uri": {
				Priority:          1,
				Type:              "string",
				Description:       "URI of the element (track / episode) to add to the playlist",
				AllowedComponents: []string{"spotify:track:id", "spotify:episode:id"},
				Required:          true,
			},
		},
		Method: addItemToPlaylist,
	}
}
