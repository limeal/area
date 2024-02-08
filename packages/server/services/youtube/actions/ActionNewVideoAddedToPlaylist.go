package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"
)

// It checks if the latest video in the playlist is the same as the one stored in the store, and if
// not, it returns the new video's information
func newVideoAddedToPlaylist(req static.AreaRequest) shared.AreaResponse {
	encode, _, err := req.Service.Endpoints["GetPlaylistItemsEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["req:playlist:id"].(string),
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	items := common.YoutubePlaylistItemsResponse{}
	if err := json.Unmarshal(encode, &items); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbItems := len(items.Items)
	if ok, errL := utils.IsLatestBasic(req.Store, nbItems); errL != nil || !ok {
		return shared.AreaResponse{Success: false}
	}

	newItem := items.Items[nbItems-1]
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"youtube:video:id":           newItem.ID,
			"youtube:video:title":        newItem.Snippet.Title,
			"youtube:video:description":  newItem.Snippet.Description,
			"youtube:video:channel":      newItem.Snippet.ChannelTitle,
			"youtube:video:published_at": newItem.Snippet.PublishedAt,
		},
	}
}

// It returns a `static.ServiceArea` object that describes the service area
// `new_video_added_to_playlist`
func DescriptorForYoutubeActionNewVideoAddedToPlaylist() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_video_added_to_playlist",
		Description: "When a new video is added to a playlist (<50)",
		Method:      newVideoAddedToPlaylist,
		RequestStore: map[string]static.StoreElement{
			"req:playlist:id": {
				Type:        "select_uri",
				Description: "Playlist where the video is added",
				Required:    true,
				Values:      []string{"/playlists"},
			},
		},
		Components: []string{
			"youtube:video:id",
			"youtube:video:title",
			"youtube:video:description",
			"youtube:video:channel",
			"youtube:video:published_at",
		},
	}
}
