package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"area-server/utils"
	"encoding/json"
)

// `NewStreamStartedByUserResponse` is a struct with a field `Data` of type `[]common.TwitchStream`.
// @property {[]common.TwitchStream} Data - An array of streams that are live.
type NewStreamStartedByUserResponse struct {
	Data []common.TwitchStream `json:"data"`
}

// It checks if a new stream has started for a user
func hasANewStreamStarted(req static.AreaRequest) shared.AreaResponse {
	userID := req.AuthStore["user_id"]

	query := make(map[string]string)
	if (*req.Store)["req:user:login"] != nil {
		if (*req.Store)["req:user:login"].(string) == "me" {
			query["user_id"] = userID.(string)
		} else {
			query["user_login"] = (*req.Store)["req:user:login"].(string)
		}
	}
	query["type"] = "live"
	query["first"] = "1"

	if (*req.Store)["ctx:game:id"] == nil && (*req.Store)["req:game:name"] != nil {
		encode, _, err := req.Service.Endpoints["GetGameEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:game:name"]})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		game := common.TwitchGames{}
		if err := json.Unmarshal(encode, &game); err != nil {
			return shared.AreaResponse{Error: err}
		}
		(*req.Store)["ctx:game:id"] = game.Data[0].ID
	}

	if (*req.Store)["ctx:game:id"] != nil {
		query["game_id"] = (*req.Store)["ctx:game:id"].(string)
	}

	encode, _, err := req.Service.Endpoints["NewStreamStartedByUserEndpoint"].CallEncode([]interface{}{req.Authorization, query})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	streams := NewStreamStartedByUserResponse{}
	if err := json.Unmarshal(encode, &streams); err != nil {
		return shared.AreaResponse{Error: err}
	}

	ok, err := utils.IsLatestByDate(req.Store, len(streams.Data), func() interface{} {
		return streams.Data[0].ID
	}, func() string {
		return streams.Data[0].StartedAt
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
			"twitch:stream:id":            streams.Data[0].ID,
			"twitch:user:id":              streams.Data[0].UserID,
			"twitch:user:login":           streams.Data[0].UserLogin,
			"twitch:user:name":            streams.Data[0].UserName,
			"twitch:game:id":              streams.Data[0].GameID,
			"twitch:game:name":            streams.Data[0].GameName,
			"twitch:stream:type":          streams.Data[0].Type,
			"twitch:stream:title":         streams.Data[0].Title,
			"twitch:stream:tag":           streams.Data[0].Tags[0],
			"twitch:stream:viewer:count":  streams.Data[0].ViewerCount,
			"twitch:stream:started:at":    streams.Data[0].StartedAt,
			"twitch:stream:language":      streams.Data[0].Language,
			"twitch:stream:thumbnail:url": streams.Data[0].ThumbnailURL,
			"twitch:stream:mature":        streams.Data[0].IsMature,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForTwitchActionNewStreamStarted() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_stream_started",
		Description: "When an user start a new stream",
		Method:      hasANewStreamStarted,
		RequestStore: map[string]static.StoreElement{
			"req:user:login": {
				Type:        "string",
				Description: "The login of the user to check (enter \"me\" for you)",
				Required:    false,
			},
			"req:game:name": {
				Priority:    1,
				Type:        "string",
				Description: "The name of the game to check (default: all)",
				Required:    false,
			},
		},
		Components: []string{
			"twitch:stream:id",
			"twitch:user:id",
			"twitch:user:login",
			"twitch:user:name",
			"twitch:game:id",
			"twitch:game:name",
			"twitch:stream:type",
			"twitch:stream:title",
			"twitch:stream:tag",
			"twitch:stream:viewer:count",
			"twitch:stream:started:at",
			"twitch:stream:language",
			"twitch:stream:thumbnail:url",
			"twitch:stream:mature",
		},
	}
}
