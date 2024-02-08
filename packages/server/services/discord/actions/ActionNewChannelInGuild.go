package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"regexp"
)

// `GuildChannelsResponse` is a struct with a field `Channels` which is an array of structs with fields
// `ID` and `Name`.
// @property {[]struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 	}} Channels - An array of channels.
type GuildChannelsResponse struct {
	Channels []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channels"`
}

// It checks if a new channel has been created in a guild
func onNewChannelInGuild(req static.AreaRequest) shared.AreaResponse {

	guildID := req.AuthStore["guild_id"]

	if (*req.Store)["req:guild:id"] != nil {
		guildID = (*req.Store)["req:guild:id"]
	}

	encode, _, err := req.Service.Endpoints["ListGuildChannelsEndpoint"].CallEncode([]interface{}{
		guildID,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	guildChannels := GuildChannelsResponse{}
	errr := json.Unmarshal(encode, &guildChannels.Channels)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbChannels := len(guildChannels.Channels)
	last := nbChannels - 1
	ok, errL := utils.IsLatestBasic(req.Store, nbChannels)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:channel:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:channel:regex"].(string), guildChannels.Channels[last].Name)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New channel in Guild (Guild ID: "+guildID.(string)+") (Channel ID:"+guildChannels.Channels[last].ID+") (Channel Name:"+guildChannels.Channels[last].Name+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:guild:id":     guildID,
			"discord:channel:id":   guildChannels.Channels[last].ID,
			"discord:channel:name": guildChannels.Channels[last].Name,
		},
	}
}

// It returns a static.ServiceArea{} struct that describes the service area
func DescriptorForDiscordActionNewChannelInGuild() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_channel_in_guild",
		Description: "When a new channel is created in a guild",
		RequestStore: map[string]static.StoreElement{
			"req:channel:regex": {
				Priority:    1,
				Type:        "string",
				Description: "The regex to match the channel name",
				Required:    false,
			},
			"req:guild:id": {
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: onNewChannelInGuild,
		Components: []string{
			"discord:guild:id",
			"discord:channel:id",
			"discord:channel:name",
		},
	}
}
