package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"area-server/utils"
	"encoding/json"
	"errors"
)

// It checks if the user is following a new streamer
func isFollowingNewStreamer(req static.AreaRequest) shared.AreaResponse {
	userID := req.AuthStore["user_id"]

	query := make(map[string]string)
	query["from_id"] = userID.(string)

	if (*req.Store)["ctx:streamer:id"] == nil && ((*req.Store)["req:streamer:login"] != nil && (*req.Store)["req:streamer:login"] != "") {
		encode, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:streamer:login"]})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		streamer := common.TwitchUsers{}
		if err := json.Unmarshal(encode, &streamer); err != nil {
			return shared.AreaResponse{Error: err}
		}

		if len(streamer.Data) == 0 {
			return shared.AreaResponse{Error: errors.New("Streamer not found")}
		}

		(*req.Store)["ctx:streamer:id"] = streamer.Data[0].ID
	}

	if (*req.Store)["ctx:streamer:id"] != nil {
		query["to_id"] = (*req.Store)["ctx:streamer:id"].(string)
	}

	body, _, err := req.Service.Endpoints["GetUserFollowsEndpoint"].CallEncode([]interface{}{req.Authorization, query})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	follows := common.TwitchFollows{}
	if err := json.Unmarshal(body, &follows); err != nil {
		return shared.AreaResponse{Error: err}
	}

	ok, errL := utils.IsLatestByDate(req.Store, len(follows.Data), func() interface{} {
		return follows.Data[0].ToId
	}, func() string {
		return follows.Data[0].FollowedAt
	})

	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"twitch:user:id":    follows.Data[0].ToId,
			"twitch:user:login": follows.Data[0].ToLogin,
			"twitch:user:name":  follows.Data[0].ToName,
			"twitch:user:date":  follows.Data[0].FollowedAt,
		},
	}
}

// `DescriptorForTwitchActionFollowNewStreamer` returns a `static.ServiceArea` that describes the
// action of following a new streamer
func DescriptorForTwitchActionFollowNewStreamer() static.ServiceArea {
	return static.ServiceArea{
		Name:        "follow_new_streamer",
		Description: "When you follow a new streamer",
		RequestStore: map[string]static.StoreElement{
			"req:streamer:login": {
				Description: "If provided, the action will be triggered only if the streamer is the one provided",
				Required:    false,
			},
		},
		Method: isFollowingNewStreamer,
		Components: []string{
			"twitch:user:id",
			"twitch:user:login",
			"twitch:user:name",
			"twitch:user:date",
		},
	}
}
