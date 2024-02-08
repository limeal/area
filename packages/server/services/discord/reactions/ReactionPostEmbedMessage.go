package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It posts an embed message to a channel
func postEmbedMessage(req static.AreaRequest) shared.AreaResponse {

	title := utils.GenerateFinalComponent((*req.Store)["req:embed:title"].(string), req.ExternalData, []string{})
	description := utils.GenerateFinalComponent((*req.Store)["req:embed:description"].(string), req.ExternalData, []string{})

	channelID := utils.GenerateFinalComponent((*req.Store)["req:channel:id"].(string), req.ExternalData, []string{
		"discord:channel:id",
	})

	if len(description) > 4096 {
		description = description[:4096]
	}

	if len(title) > 256 {
		title = title[:256]
	}

	body := make(map[string]interface{})
	body["embeds"] = []map[string]interface{}{
		{
			"title":       title,
			"description": description,
		},
	}

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

	req.Logger.WriteInfo("[Reaction] Post embed message (Channel ID: "+channelID+") (Title: "+title+") (Description: "+description+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForDiscordReactionPostEmbedMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "post_embed_message",
		Description: "Post an embed message to a channel",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:              "select_uri",
				Values:            []string{"/guilds/${req:guild:id}/channels"},
				Description:       "The ID of the channel to send the message to",
				Required:          true,
				AllowedComponents: []string{"discord:channel:id"},
			},
			"req:embed:title": {
				Priority:    1,
				Type:        "string",
				Description: "The content of the message",
				Required:    true,
			},
			"req:embed:description": {
				Priority:    2,
				Type:        "long_string",
				Description: "The content of the message",
				Required:    true,
			},
			"req:not:reply": {
				Priority:    3,
				Type:        "select",
				Description: "Whether to reply to the message (useless if action is not a discord message, default: true)",
				Values:      []string{"true", "false"},
				Required:    false,
			},
			"req:guild:id": {
				Priority:    3,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: postEmbedMessage,
	}
}
