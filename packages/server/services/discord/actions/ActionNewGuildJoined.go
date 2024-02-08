package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
	"errors"
	"regexp"
)

// `PartialGuild` is a struct with fields `ID`, `Name`, `Owner`, `Icon`, and `Features`.
//
// The `json` tags are used to tell the JSON encoder/decoder how to map the fields to JSON.
//
// The `json:"id"` tag tells the JSON encoder/decoder that the `ID` field should be mapped to the `id`
// key in the JSON.
//
// The `json:"name"` tag tells the JSON encoder/decoder that the `Name` field should be mapped to the `
// @property {string} ID - The ID of the guild.
// @property {string} Name - The name of the guild.
// @property {bool} Owner - Whether the user is the owner of the guild.
// @property {string} Icon - The hash of the guild's icon.
// @property {[]string} Features - A list of guild features.
type PartialGuild struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Owner    bool     `json:"owner"`
	Icon     string   `json:"icon"`
	Features []string `json:"features"`
}

// `UserGuildsResponse` is a struct with a field `Guilds` of type `[]PartialGuild`.
//
// The `[]` means that `Guilds` is a slice of `PartialGuild`s.
//
// The `PartialGuild` type is defined in the same file as `UserGuildsResponse`, so we can just look at
// it:
// @property {[]PartialGuild} Guilds - An array of PartialGuilds.
type UserGuildsResponse struct {
	Guilds []PartialGuild `json:"guilds"`
}

// It checks if the user has joined a new guild, and if so, it returns the guild's information
func onNewGuildJoined(req static.AreaRequest) shared.AreaResponse {

	encode, _, err := req.Service.Endpoints["GetUserGuildsEndpoint"].CallEncode([]interface{}{req.Authorization})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	userGuilds := UserGuildsResponse{}
	errr := json.Unmarshal(encode, &userGuilds.Guilds)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbGuilds := len(userGuilds.Guilds)
	if (*req.Store)["ctx:guilds:number"] == nil {
		(*req.Store)["ctx:guilds:number"] = nbGuilds
		(*req.Store)["ctx:guilds"] = userGuilds.Guilds
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:guilds:number"] == nbGuilds {
		return shared.AreaResponse{Success: false}
	}

	if nbGuilds < (*req.Store)["ctx:guilds:number"].(int) {
		(*req.Store)["ctx:guilds:number"] = nbGuilds
		(*req.Store)["ctx:guilds"] = userGuilds.Guilds
		return shared.AreaResponse{Success: false}
	}

	if nbGuilds <= 0 {
		return shared.AreaResponse{Error: errors.New("Invalid number of guilds")}
	}

	// Find the new guild in the list
	newGuild := userGuilds.Guilds[0]
	for _, guild := range userGuilds.Guilds {
		found := false
		for _, oldGuild := range (*req.Store)["ctx:guilds"].([]PartialGuild) {
			if guild.ID == oldGuild.ID {
				found = true
				break
			}
		}
		if !found {
			newGuild = guild
			break
		}
	}

	if (*req.Store)["req:guild:name"] != nil && (*req.Store)["req:guild:name"] != newGuild.Name {
		(*req.Store)["ctx:guilds:number"] = nbGuilds
		(*req.Store)["ctx:guilds"] = userGuilds.Guilds
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:guild:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:guild:regex"].(string), newGuild.Name)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			(*req.Store)["ctx:guilds:number"] = nbGuilds
			(*req.Store)["ctx:guilds"] = userGuilds.Guilds
			return shared.AreaResponse{Success: false}
		}
	}

	(*req.Store)["ctx:guilds:number"] = nbGuilds
	(*req.Store)["ctx:guilds"] = userGuilds.Guilds
	req.Logger.WriteInfo("[Action] New guild joined (Name: "+newGuild.Name+") (ID: "+newGuild.ID+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:guild:id":       newGuild.ID,
			"discord:guild:name":     newGuild.Name,
			"discord:guild:icon":     newGuild.Icon,
			"discord:guild:owner":    newGuild.Owner,
			"discord:guild:features": newGuild.Features,
		},
	}
}

// `DescriptorForDiscordActionNewGuildJoined` returns a `static.ServiceArea` that describes the
// `onNewGuildJoined` function
func DescriptorForDiscordActionNewGuildJoined() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_guild_joined",
		Description: "When a new guild is joined",
		RequestStore: map[string]static.StoreElement{
			"req:guild:name": {
				Type:        "string",
				Description: "The Name of the guild that check if joined in",
				Required:    false,
			},
			"req:guild:regex": {
				Priority:    1,
				Type:        "string",
				Description: "The regex of the guild that check if joined in",
				Required:    false,
			},
		},
		Method: onNewGuildJoined,
		Components: []string{
			"discord:guild:id",
			"discord:guild:name",
			"discord:guild:icon",
			"discord:guild:owner",
			"discord:guild:features",
		},
	}
}
