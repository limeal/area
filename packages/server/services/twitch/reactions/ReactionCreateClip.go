package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"encoding/json"
	"errors"
)

// It creates a clip of the streamer
func createClip(req static.AreaRequest) shared.AreaResponse {
	if (*req.Store)["ctx:broadcaster:id"] == nil {
		switch (*req.Store)["req:streamer:entity:type"] {
		case "login":
			encode, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:streamer:entity:value"]})
			if err != nil {
				return shared.AreaResponse{Error: err}
			}
			user := common.TwitchUsers{}
			if err := json.Unmarshal(encode, &user); err != nil {
				return shared.AreaResponse{Error: err}
			}

			if len(user.Data) == 0 {
				return shared.AreaResponse{Error: errors.New("Streamer not found")}
			}

			(*req.Store)["ctx:broadcaster:id"] = user.Data[0].ID
		case "id":
			(*req.Store)["ctx:broadcaster:id"] = (*req.Store)["req:streamer:entity:value"]
		default:
			return shared.AreaResponse{Error: errors.New("Invalid streamer type")}
		}
	}

	_, httpResp, err := req.Service.Endpoints["CreateClipEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["ctx:broadcaster:id"].(string)})

	if httpResp != nil && httpResp.StatusCode == 404 {
		req.Logger.WriteInfo("(Twitch) No clip created, the streamer is not living", true)
		return shared.AreaResponse{
			Success: false,
		}
	}

	return shared.AreaResponse{Error: err}
}

// It returns a `static.ServiceArea` struct, which contains the name, description, request store and
// the method to be executed
func DescriptorForTwitchReactionCreateClip() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_clip",
		Description: "Create a new clip",
		RequestStore: map[string]static.StoreElement{
			"req:streamer:entity:type": {
				Description: "The type of the streamer, the clip will be created for",
				Type:        "select",
				Required:    true,
				Values: []string{
					"login",
					"id",
				},
			},
			"req:streamer:entity:value": {
				Description: "The login or id of the streamer, the clip will be created for",
				Required:    true,
			},
		},
		Method: createClip,
	}
}
