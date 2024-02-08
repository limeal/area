package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
	"errors"

	"github.com/forPelevin/gomoji"
)

// `UsersThatReactedResponse` is a struct that contains a slice of structs that contain three strings.
// @property {[]struct {
// 		ID            string `json:"id"`
// 		Username      string `json:"username"`
// 		Discriminator string `json:"discriminator"`
// 	}} Users - An array of users that reacted to the message.
type UsersThatReactedResponse struct {
	Users []struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
	}
}

// It gets the last user that reacted to a message with a specific emoji
func onNewReactionAddedToMessage(req static.AreaRequest) shared.AreaResponse {

	emoji := gomoji.FindAll((*req.Store)["req:emoji"].(string))
	if len(emoji) != 1 {
		return shared.AreaResponse{Error: errors.New("Invalid emoji")}
	}

	query := make(map[string]string)
	query["limit"] = "1"

	if (*req.Store)["ctx:last:react"] != nil {
		query["after"] = (*req.Store)["ctx:last:react"].(string)
	}

	encode, _, err := req.Service.Endpoints["ListUserThatReactedEndpoint"].CallEncode([]interface{}{
		(*req.Store)["req:channel:id"],
		(*req.Store)["req:message:id"],
		emoji[0].Character,
		query,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	usersThatReacted := UsersThatReactedResponse{}
	errr := json.Unmarshal(encode, &usersThatReacted.Users)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbReacts := len(usersThatReacted.Users)
	if (*req.Store)["ctx:reacts:number"] == nil {
		(*req.Store)["ctx:reacts:number"] = nbReacts
	}

	if (*req.Store)["ctx:reacts:number"].(int) > nbReacts {
		(*req.Store)["ctx:reacts:number"] = nbReacts
	}

	if (*req.Store)["ctx:reacts:number"].(int) == 0 && nbReacts == 0 {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:reacts:number"].(int) > 0 && (*req.Store)["ctx:last:react"] == nil {
		(*req.Store)["ctx:last:react"] = usersThatReacted.Users[0].ID
		return shared.AreaResponse{Success: false}
	}

	if nbReacts <= 0 {
		return shared.AreaResponse{Error: errors.New("Invalid number of reacts")}
	}

	(*req.Store)["ctx:reacts:number"] = nbReacts
	(*req.Store)["ctx:last:react"] = usersThatReacted.Users[0].ID
	req.Logger.WriteInfo("[Action] New reaction added by user (Channel ID: "+(*req.Store)["req:channel:id"].(string)+") (Message: "+(*req.Store)["req:message:id"].(string)+") (Emoji: "+emoji[0].Character+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:user:id":            usersThatReacted.Users[0].ID,
			"discord:user:username":      usersThatReacted.Users[0].Username,
			"discord:user:discriminator": usersThatReacted.Users[0].Discriminator,
			"discord:emoji":              emoji[0].Character,
		},
	}
}

// It returns a static.ServiceArea{} that describes the service area
func DescriptorForDiscordActionNewReactionAddedByUser() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_reaction_added_by_user",
		Description: "When a new reaction is added to a message by an user",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "select_uri",
				Values:      []string{"/guilds/${req:guild:id}/channels"},
				Description: "The ID of the channel to send the message to",
				Required:    true,
			},
			"req:message:id": {
				Priority:    1,
				Type:        "select_uri",
				Values:      []string{"/channels/${req:channel:id}/messages?limit=10"},
				NeedFields:  []string{"req:channel:id"},
				Description: "The ID of the message to add the reaction to",
				Required:    true,
			},
			"req:emoji": {
				Priority:    2,
				Type:        "string",
				Description: "The emoji to add as a reaction",
				Required:    true,
			},
			"req:guild:id": {
				Priority:    3,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: onNewReactionAddedToMessage,
		Components: []string{
			"discord:user:id",
			"discord:user:username",
			"discord:user:discriminator",
			"discord:emoji",
		},
	}
}
