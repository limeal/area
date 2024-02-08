package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It posts a message to a channel
func postMessage(req static.AreaRequest) shared.AreaResponse {

	content := utils.GenerateFinalComponent((*req.Store)["req:message:content"].(string), req.ExternalData, []string{})

	body := make(map[string]interface{})
	body["content"] = content

	channelID := utils.GenerateFinalComponent((*req.Store)["req:channel:id"].(string), req.ExternalData, []string{
		"discord:channel:id",
	})

	if req.ExternalData["discord:message:id"] != nil && (*req.Store)["req:not:reply"] != "true" {

		_, _, err0 := req.Service.Endpoints["GetChannelMessageEndpoint"].CallEncode([]interface{}{
			channelID,
			req.ExternalData["discord:message:id"],
		})

		if err0 == nil {
			body["message_reference"] = map[string]interface{}{
				"message_id": req.ExternalData["discord:message:id"],
			}
		}
	}

	fbody, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["CreateMessageEndpoint"].CallEncode([]interface{}{
		channelID,
		string(fbody),
	})
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	req.Logger.WriteInfo("[Reaction] Post message (Channel ID: "+channelID+") (Content: "+content+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForDiscordReactionPostMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "post_message",
		Description: "Post a message to a channel",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "select_uri",
				Values:      []string{"/guilds/${req:guild:id}/channels"},
				Description: "The ID of the channel to send the message to",
				Required:    true,
			},
			"req:message:content": {
				Priority:    1,
				Type:        "long_string",
				Description: "The content of the message",
				Required:    true,
			},
			"req:not:reply": {
				Priority:    2,
				Type:        "select",
				Description: "Whether to reply to the message (useless if action is not a discord message, default: true)",
				Values:      []string{"true", "false"},
				Required:    false,
			},
			"req:guild:id": {
				Priority:    2,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: postMessage,
	}
}
