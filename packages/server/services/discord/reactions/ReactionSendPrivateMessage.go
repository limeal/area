package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// A Message is a struct with three fields: ID, ChannelID, and Content.
// @property {string} ID - The unique ID of the message.
// @property {string} ChannelID - The ID of the channel that the message belongs to.
// @property {string} Content - The message content
type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

// It sends a private message to a user
func sendPrivateMessage(req static.AreaRequest) shared.AreaResponse {

	userID := utils.GenerateFinalComponent((*req.Store)["req:user:id"].(string), req.ExternalData, []string{
		"discord:user:id",
	})

	body, _, err0 := req.Service.Endpoints["CreateDMEndpoint"].Call([]interface{}{
		userID,
	})
	if err0 != nil {
		return shared.AreaResponse{Error: err0}
	}

	fbody, err1 := json.Marshal(map[string]interface{}{
		"content": utils.GenerateFinalComponent((*req.Store)["req:message:content"].(string), req.ExternalData, []string{}),
	})
	if err1 != nil {
		return shared.AreaResponse{Error: err1}
	}

	emessage, _, err2 := req.Service.Endpoints["CreateMessageEndpoint"].CallEncode([]interface{}{
		body["id"].(string),
		string(fbody),
	})
	if err2 != nil {
		return shared.AreaResponse{Error: err2}
	}

	message := Message{}
	err3 := json.Unmarshal(emessage, &message)
	if err3 != nil {
		return shared.AreaResponse{Error: err3}
	}

	req.Logger.WriteInfo("[Reaction] Sent private message to user "+userID+" with content "+message.Content+" (Message ID: "+message.ID+")"+" (Channel ID: "+message.ChannelID+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForDiscordReactionSendPrivateMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "send_private_message",
		Description: "Send a private message to a channel",
		RequestStore: map[string]static.StoreElement{
			"req:user:id": {
				Type:        "string",
				Description: "The ID of the user to send the message to",
				Required:    true,
			},
			"req:message:content": {
				Priority:    1,
				Type:        "long_string",
				Description: "The content of the message",
				Required:    true,
			},
		},
		Method: sendPrivateMessage,
	}
}
