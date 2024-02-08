package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"area-server/utils"
	"encoding/json"
	"time"
)

// It checks if a new clip has been captured on Twitch
func hasANewClipBeenCaptured(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)

	if (*req.Store)["ctx:clip:created_at"] == nil {
		(*req.Store)["ctx:clip:created_at"] = time.Now().Format(time.RFC3339)
	}

	query["first"] = "1"
	query["started_at"] = (*req.Store)["ctx:clip:created_at"].(string)

	if (*req.Store)["ctx:entity:id"] == nil {
		switch (*req.Store)["req:clip:entity:type"] {
		case "user":
			encode, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:clip:entity:value"]})
			if err != nil {
				return shared.AreaResponse{Error: err}
			}
			streamer := common.TwitchUsers{}
			if err := json.Unmarshal(encode, &streamer); err != nil {
				return shared.AreaResponse{Error: err}
			}
			(*req.Store)["ctx:entity:id"] = streamer.Data[0].ID
		case "game":
			encode, _, err := req.Service.Endpoints["GetGameEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:clip:entity:value"]})
			if err != nil {
				return shared.AreaResponse{Error: err}
			}
			game := common.TwitchStreams{}
			if err := json.Unmarshal(encode, &game); err != nil {
				return shared.AreaResponse{Error: err}
			}
			(*req.Store)["ctx:entity:id"] = game.Data[0].ID
		}
	}

	switch (*req.Store)["req:clip:entity:type"] {
	case "user":
		query["broadcaster_id"] = (*req.Store)["ctx:entity:id"].(string)
	case "game":
		query["game_id"] = (*req.Store)["ctx:entity:id"].(string)
	}

	encode, _, err := req.Service.Endpoints["GetClipsEndpoint"].CallEncode([]interface{}{req.Authorization, query})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	clips := common.TwitchClips{}
	if err := json.Unmarshal(encode, &clips); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbClips := len(clips.Data)
	ok, err := utils.IsLatestByDate(req.Store, nbClips, func() interface{} {
		return clips.Data[0].ID
	}, func() string {
		return clips.Data[0].CreatedAt
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"twitch:clip:id":            clips.Data[0].ID,
			"twitch:clip:url":           clips.Data[0].URL,
			"twitch:clip:embed_url":     clips.Data[0].EmbedURL,
			"twitch:broadcaster:id":     clips.Data[0].BroadcasterID,
			"twitch:broadcaster:name":   clips.Data[0].BroadcasterName,
			"twitch:creator:id":         clips.Data[0].CreatorID,
			"twitch:creator:name":       clips.Data[0].CreatorName,
			"twitch:video:id":           clips.Data[0].VideoID,
			"twitch:game:id":            clips.Data[0].GameID,
			"twitch:clip:language":      clips.Data[0].Language,
			"twitch:clip:title":         clips.Data[0].Title,
			"twitch:clip:view_count":    clips.Data[0].ViewCount,
			"twitch:clip:created_at":    clips.Data[0].CreatedAt,
			"twitch:clip:thumbnail_url": clips.Data[0].ThumbnailURL,
			"twitch:clip:duration":      clips.Data[0].Duration,
		},
	}
}

// It returns a static.ServiceArea struct that describes the action
func DescriptorForTwitchActionNewClipCaptured() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_clip_captured",
		Description: "Triggered when a new clip is captured",
		RequestStore: map[string]static.StoreElement{
			"req:clip:entity:value": {
				Priority:    1,
				Description: "Can be either a user login or name of a game",
				Required:    true,
			},
			"req:clip:entity:type": {
				Description: "Choose between user or game",
				Required:    true,
				Type:        "select",
				Values: []string{
					"user",
					"game",
				},
			},
		},
		Method: hasANewClipBeenCaptured,
		Components: []string{
			"twitch:clip:id",
			"twitch:clip:url",
			"twitch:clip:embed_url",
			"twitch:broadcaster:id",
			"twitch:broadcaster:name",
			"twitch:creator:id",
			"twitch:creator:name",
			"twitch:video:id",
			"twitch:game:id",
			"twitch:clip:title",
			"twitch:clip:language",
			"twitch:clip:view_count",
			"twitch:clip:created_at",
			"twitch:clip:thumbnail_url",
			"twitch:clip:duration",
		},
	}
}
