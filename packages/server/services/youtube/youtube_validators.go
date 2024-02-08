package youtube

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/youtube/common"
	"encoding/json"
)

// It returns a map of validators that can be used to validate the request parameters of the Youtube
// service
func YoutubeValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:video:id":                VideoIdValidator,
		"req:channel:id":              ChannelIdValidator,
		"req:playlist:id":             PlaylistIdValidator,
		"req:entity:commentable:type": EntityCommentableTypeValidator,
		"req:entity:commentable:id":   EntityCommentableIdValidator,
	}
}

// ------------------ Validators ------------------

// It checks if the value is a valid YouTube video ID
func VideoIdValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{youtube:video:id}}" {
		return true
	}

	encode, _, err := service.Endpoints["GetAllVideosEndpoint"].CallEncode([]interface{}{
		authorization,
		map[string]string{
			"id":   value.(string),
			"part": "id",
		},
	})

	if err != nil {
		return false
	}

	var videos common.YoutubeVideoResponse
	if err := json.Unmarshal(encode, &videos); err != nil {
		return false
	}

	if len(videos.Items) == 0 {
		return false
	}

	return true
}

// It checks if the value is a valid YouTube channel ID
func ChannelIdValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{youtube:channel:id}}" {
		return true
	}

	encode, _, err := service.Endpoints["GetAllChannelsEndpoint"].CallEncode([]interface{}{
		authorization,
		map[string]string{
			"id":   value.(string),
			"part": "id",
		},
	})

	if err != nil {
		return false
	}

	var channels common.YoutubeChannelsResponse
	if err := json.Unmarshal(encode, &channels); err != nil {
		return false
	}

	if len(channels.Items) == 0 {
		return false
	}

	return true
}

// It checks if the playlist id is valid
func PlaylistIdValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	encode, _, err := service.Endpoints["GetAllPlaylistsEndpoint"].CallEncode([]interface{}{
		authorization,
		map[string]string{
			"id":   value.(string),
			"part": "id",
		},
	})

	if err != nil {
		return false
	}

	var playlists common.YoutubePlaylistsResponse
	if err := json.Unmarshal(encode, &playlists); err != nil {
		return false
	}

	if len(playlists.Items) == 0 {
		return false
	}

	return true
}

// It takes a string, and checks if it's a valid YouTube channel name
func ChannelNameValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	encode, _, err := service.Endpoints["GetAllChannelsEndpoint"].CallEncode([]interface{}{
		authorization,
		map[string]string{
			"forUsername": value.(string),
			"part":        "id",
		},
	})

	if err != nil {
		return false
	}

	var channels common.YoutubeChannelsResponse
	if err := json.Unmarshal(encode, &channels); err != nil {
		return false
	}

	if len(channels.Items) == 0 {
		return false
	}

	return true
}

// It checks if the commentable id is valid
func EntityCommentableIdValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:entity:commentable:type"] == nil {
		return false
	}

	switch store["req:entity:commentable:type"].(string) {
	case "video":
		if VideoIdValidator(authorization, service, value, store) {
			return true
		}
	case "channel":
		if ChannelIdValidator(authorization, service, value, store) {
			return true
		}
	case "comment":
		if value == "{{youtube:comment:id}}" {
			return true
		}

		encode, _, err := service.Endpoints["GetAllRepliesEndpoint"].CallEncode([]interface{}{
			authorization,
			map[string]string{
				"id":   value.(string),
				"part": "id",
			},
		})

		if err != nil {
			return false
		}

		var comments common.YoutubeCommentsResponse
		if err := json.Unmarshal(encode, &comments); err != nil {
			return false
		}

		if len(comments.Items) == 0 {
			return false
		}

		return true
	}

	return false
}

// It checks that the value is a string and that it's either "video", "channel" or "comment"
func EntityCommentableTypeValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	switch value.(string) {
	case "video", "channel", "comment":
		return true
	}

	return false
}
