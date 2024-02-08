package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"area-server/utils"
	"encoding/json"
	"errors"
)

// It sends a whisper to a user
func sendWhisper(req static.AreaRequest) shared.AreaResponse {
	userID := req.AuthStore["user_id"].(string)

	message := utils.GenerateFinalComponent((*req.Store)["req:message:content"].(string), req.ExternalData, []string{})

	if (*req.Store)["ctx:receiver:id"] == nil {
		encode, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:user:login"]})
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

		(*req.Store)["ctx:receiver:id"] = streamer.Data[0].ID
	}

	_, httpResp, err := req.Service.Endpoints["SendWhisperEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userID,
		(*req.Store)["ctx:receiver:id"].(string),
		message,
	})

	if httpResp != nil && httpResp.StatusCode == 400 {
		return shared.AreaResponse{Error: errors.New("User not allow whispers")}
	}

	if httpResp != nil && httpResp.StatusCode == 401 {
		return shared.AreaResponse{Error: errors.New("You don't have a verified phone number")}
	}

	return shared.AreaResponse{
		Error: err,
	}
}

// `DescriptorForTwitchReactionSendWhisper` returns a `static.ServiceArea` that describes the
// `send_whisper` service area
func DescriptorForTwitchReactionSendWhisper() static.ServiceArea {
	return static.ServiceArea{
		Name:        "send_whisper",
		Description: "Send a whisper to a user",
		RequestStore: map[string]static.StoreElement{
			"req:user:login": {
				Description: "The login of the user to send the whisper to",
				Required:    true,
			},
			"req:message:content": {
				Priority:    1,
				Description: "The message to send",
				Required:    true,
			},
		},
		Method: sendWhisper,
	}
}
