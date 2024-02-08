package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/discord/common"
	"area-server/utils"
	"encoding/json"
)

// It gets the latest pinned message in a channel, and returns it
func onNewPinnedMessageInChannel(req static.AreaRequest) shared.AreaResponse {

	if common.IsRateLimited(req.Store) {
		return shared.AreaResponse{Success: false}
	}

	encode, httpResp, err := req.Service.Endpoints["ListPinnedMessagesEndpoint"].CallEncode([]interface{}{
		(*req.Store)["req:channel:id"],
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if ok, err := common.SetRateLimit(httpResp, req.Store, encode); !ok || err != nil {
		return shared.AreaResponse{Error: err, Success: false}
	}

	pinnedMessages := ChannelMessagesResponse{}
	errr := json.Unmarshal(encode, &pinnedMessages.Messages)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbPinnedMessages := len(pinnedMessages.Messages)
	ok, errL := utils.IsLatestBasic(req.Store, nbPinnedMessages)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:allow:bot"] != nil && (*req.Store)["req:allow:bot"].(string) == "false" {
		if pinnedMessages.Messages[0].Author.Bot {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New pinned message in channel (Channel: + "+(*req.Store)["req:channel:id"].(string)+")"+" (Message: "+pinnedMessages.Messages[0].Content+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:channel:id":         (*req.Store)["req:channel:id"].(string),
			"discord:message:id":         pinnedMessages.Messages[0].ID,
			"discord:message:content":    pinnedMessages.Messages[0].Content,
			"discord:user:id":            pinnedMessages.Messages[0].Author.ID,
			"discord:user:username":      pinnedMessages.Messages[0].Author.Username,
			"discord:user:discriminator": pinnedMessages.Messages[0].Author.Discriminator,
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForDiscordActionNewPinnedMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_pinned_message",
		Description: "When a new pinned message is created",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "select_uri",
				Values:      []string{"/guilds/${req:guild:id}/channels"},
				Description: "The ID of the channel to send the message to",
				Required:    true,
			},
			"req:guild:id": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
			"req:allow:bot": {
				Priority:    1,
				Type:        "select",
				Description: "Allow bot to trigger this action",
				Required:    false,
				Values:      []string{"true", "false"},
			},
		},
		Method: onNewPinnedMessageInChannel,
		Components: []string{
			"discord:channel:id",
			"discord:message:id",
			"discord:message:content",
			"discord:user:id",
			"discord:user:username",
			"discord:user:discriminator",
		},
	}
}
