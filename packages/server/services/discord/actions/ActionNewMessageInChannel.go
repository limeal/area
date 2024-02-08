package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/discord/common"
	"encoding/json"
	"errors"
	"regexp"
)

// `ChannelMessagesResponse` is a struct with a field `Messages` which is a slice of structs with
// fields `ID`, `ChannelID`, `Author`, and `Content`.
//
// The `Author` field is a struct with fields `ID`, `Username`, `Discriminator`, and `Bot`.
//
// The `Author` field is a struct with fields `ID`, `Username`, `Discriminator`, and `Bot`.
//
// The `Author` field is a struct with fields `ID`, `Username`, `Discriminator`, and
// @property {[]struct {
// 		ID        string `json:"id"`
// 		ChannelID string `json:"channel_id"`
// 		Author    struct {
// 			ID            string `json:"id"`
// 			Username      string `json:"username"`
// 			Discriminator string `json:"discriminator"`
// 			Bot           bool   `json:"bot"`
// 		} `json:"author"`
// 		Content string `json:"content"`
// 	}} Messages - An array of messages.
type ChannelMessagesResponse struct {
	Messages []struct {
		ID        string `json:"id"`
		ChannelID string `json:"channel_id"`
		Author    struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Discriminator string `json:"discriminator"`
			Bot           bool   `json:"bot"`
		} `json:"author"`
		Content string `json:"content"`
	}
}

// It checks if a new message has been posted in a channel
func onNewMessageInChannel(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)

	query["limit"] = "1"
	if (*req.Store)["ctx:last:message"] != nil {
		query["after"] = (*req.Store)["ctx:last:message"].(string)
	}

	if common.IsRateLimited(req.Store) {
		return shared.AreaResponse{Success: false}
	}

	encode, httpResp, err := req.Service.Endpoints["ListChannelMessagesEndpoint"].CallEncode([]interface{}{
		(*req.Store)["req:channel:id"],
		query,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if ok, err := common.SetRateLimit(httpResp, req.Store, encode); !ok || err != nil {
		return shared.AreaResponse{Error: err, Success: false}
	}

	messageList := ChannelMessagesResponse{}
	errr := json.Unmarshal(encode, &messageList.Messages)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbMessages := len(messageList.Messages)

	if (*req.Store)["ctx:messages:number"] == nil {
		(*req.Store)["ctx:messages:number"] = nbMessages
	}

	if (*req.Store)["ctx:messages:number"].(int) > nbMessages {
		(*req.Store)["ctx:messages:number"] = nbMessages
	}

	if (*req.Store)["ctx:messages:number"].(int) == 0 && nbMessages == 0 {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:messages:number"].(int) > 0 && (*req.Store)["ctx:last:message"] == nil {
		(*req.Store)["ctx:last:message"] = messageList.Messages[0].ID
		return shared.AreaResponse{Success: false}
	}

	if nbMessages <= 0 {
		return shared.AreaResponse{Error: errors.New("Invalid number of messages")}
	}

	if (*req.Store)["req:allow:bot"] != nil && (*req.Store)["req:allow:bot"].(string) == "false" && messageList.Messages[0].Author.Bot {
		(*req.Store)["ctx:messages:number"] = nbMessages
		(*req.Store)["ctx:last:message"] = messageList.Messages[0].ID
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:user:id"] != nil && (*req.Store)["req:user:id"].(string) != messageList.Messages[0].Author.ID {
		(*req.Store)["ctx:messages:number"] = nbMessages
		(*req.Store)["ctx:last:message"] = messageList.Messages[0].ID
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:message:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:message:regex"].(string), messageList.Messages[0].Content)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			(*req.Store)["ctx:messages:number"] = nbMessages
			(*req.Store)["ctx:last:message"] = messageList.Messages[0].ID
			return shared.AreaResponse{Success: false}
		}
	}

	(*req.Store)["ctx:messages:number"] = nbMessages
	(*req.Store)["ctx:last:message"] = messageList.Messages[0].ID
	req.Logger.WriteInfo("[Action] New message in channel: (Content: "+messageList.Messages[0].Content+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:channel:id":         (*req.Store)["req:channel:id"].(string),
			"discord:message:id":         messageList.Messages[0].ID,
			"discord:message:content":    messageList.Messages[0].Content,
			"discord:user:id":            messageList.Messages[0].Author.ID,
			"discord:user:username":      messageList.Messages[0].Author.Username,
			"discord:user:discriminator": messageList.Messages[0].Author.Discriminator,
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForDiscordActionNewMessageInChannel() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_message_in_channel",
		Description: "When a new message is posted in a channel",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "select_uri",
				Values:      []string{"/guilds/${req:guild:id}/channels"},
				Description: "The channel to send the message to",
				Required:    true,
			},
			"req:allow:bot": {
				Priority:    4,
				Type:        "select",
				Description: "Allow bot messages to trigger the reaction (default: false)",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:message:regex": {
				Priority:    3,
				Type:        "string",
				Description: "A regex to match the message content against",
				Required:    false,
			},
			"req:guild:id": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
			"req:user:id": {
				Priority:    2,
				Type:        "select_uri",
				Description: "The user to check in the new channel (default: all users)",
				Required:    false,
				Values:      []string{"/guilds/${req:guild:id}/members"},
			},
		},
		Method: onNewMessageInChannel,
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
