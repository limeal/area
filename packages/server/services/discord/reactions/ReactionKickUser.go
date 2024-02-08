package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
)

// It kicks a user from a guild
func kickUserFromGuild(req static.AreaRequest) shared.AreaResponse {

	guildID := req.AuthStore["guild_id"].(string)

	if (*req.Store)["req:guild:id"] != nil {
		guildID = (*req.Store)["req:guild:id"].(string)
	}

	userID := utils.GenerateFinalComponent((*req.Store)["req:user:id"].(string), req.ExternalData, []string{"discord:user:id"})

	_, _, err := req.Service.Endpoints["KickUserEndpoint"].CallEncode([]interface{}{
		guildID,
		userID,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	req.Logger.WriteInfo("[Reaction] Kick user from guild (Guild ID: "+guildID+") (User ID: "+userID+")", true)
	return shared.AreaResponse{
		Error: nil,
	}
}

// "DescriptorForDiscordReactionKickUser returns a static.ServiceArea that describes the kick_user
// service area."
//
// The first thing we do is create a static.ServiceArea struct. This struct is used to describe the
// service area
func DescriptorForDiscordReactionKickUser() static.ServiceArea {
	return static.ServiceArea{
		Name:        "kick_user",
		Description: "Kick a user from a guild",
		RequestStore: map[string]static.StoreElement{
			"req:user:id": {
				Type:              "select_uri",
				Values:            []string{"/guilds/default/members"},
				Description:       "The ID of the user to kick",
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
		Method: kickUserFromGuild,
	}
}
