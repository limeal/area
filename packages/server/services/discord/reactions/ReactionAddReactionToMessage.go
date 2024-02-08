package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"errors"

	"github.com/forPelevin/gomoji"
)

// It adds a reaction to a message
func addReactionToMessage(req static.AreaRequest) shared.AreaResponse {

	channelID := utils.GenerateFinalComponent((*req.Store)["req:channel:id"].(string), req.ExternalData, []string{"discord:channel:id"})
	messageID := utils.GenerateFinalComponent((*req.Store)["req:message:id"].(string), req.ExternalData, []string{"discord:message:id"})
	emoji := utils.GenerateFinalComponent((*req.Store)["req:emoji"].(string), req.ExternalData, []string{"discord:emoji"})

	emojis := gomoji.FindAll(emoji)
	if len(emojis) != 1 {
		return shared.AreaResponse{Error: errors.New("Invalid emoji")}
	}

	_, _, err := req.Service.Endpoints["CreateReactionEndpoint"].CallEncode([]interface{}{
		channelID,
		messageID,
		emojis[0].Character,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	req.Logger.WriteInfo("[Reaction] Added reaction to message: (MessageID: "+messageID+") (Emoji: "+emojis[0].Character+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// "DescriptorForDiscordReactionAddReactionToMessage returns a static.ServiceArea that describes the
// add_reaction_to_message service area."
//
// The first thing you'll notice is that the function name is
// `DescriptorForDiscordReactionAddReactionToMessage`. This is the name of the service area, with the
// prefix `DescriptorFor` and the suffix `ServiceArea`
func DescriptorForDiscordReactionAddReactionToMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_reaction_to_message",
		Description: "Add a reaction to a message",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:              "select_uri",
				Values:            []string{"/guilds/default/channels"},
				Description:       "The ID of the channel to send the message to",
				Required:          true,
				AllowedComponents: []string{"discord:channel:id"},
			},
			"req:message:id": {
				Priority:          1,
				Type:              "select_uri",
				Values:            []string{"/channels/${req:channel:id}/messages"},
				NeedFields:        []string{"req:channel:id"},
				Description:       "The ID of the message to add the reaction to",
				Required:          true,
				AllowedComponents: []string{"discord:message:id"},
			},
			"req:emoji": {
				Priority:          2,
				Type:              "string",
				Description:       "The emoji to add as a reaction",
				Required:          true,
				AllowedComponents: []string{"discord:emoji"},
			},
		},
		Method: addReactionToMessage,
	}
}
