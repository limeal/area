package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
)

// It bans a user from a guild
func banUserFromGuild(req static.AreaRequest) shared.AreaResponse {

	guildID := req.AuthStore["guild_id"].(string)

	if (*req.Store)["req:guild:id"] != nil {
		guildID = (*req.Store)["req:guild:id"].(string)
	}

	userID := utils.GenerateFinalComponent((*req.Store)["req:user:id"].(string), req.ExternalData, []string{"discord:user:id"})

	_, _, err := req.Service.Endpoints["BanUserEndpoint"].CallEncode([]interface{}{guildID, userID})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	req.Logger.WriteInfo("[Reaction] Ban user from guild (Guild ID: "+guildID+") (User ID: "+userID+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// It takes a user ID and a guild ID, and bans the user from the guild
func DescriptorForDiscordReactionBanUser() static.ServiceArea {
	return static.ServiceArea{
		Name:        "ban_user",
		Description: "Ban a user from a guild",
		RequestStore: map[string]static.StoreElement{
			"req:user:id": {
				Type:              "select_uri",
				Values:            []string{"/guilds/default/members"},
				Description:       "The ID of the user to ban",
				Required:          true,
				AllowedComponents: []string{"discord:user:id"},
			},
			"req:guild:id": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: banUserFromGuild,
	}
}
