package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
	"errors"
)

// `Member` is a struct that contains a `User` struct that contains a `ID`, `Username`,
// `Discriminator`, and `Bot` field.
//
// Now that we have a type that matches the JSON, we can use the `json.Unmarshal` function to parse the
// JSON into a `Member` type.
// @property User - The user object of the member.
type Member struct {
	User struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
	} `json:"user"`
}

// `GuildMembersResponse` is a struct with a single field, `Members`, which is a slice of `Member`s.
// @property {[]Member} Members - An array of Member objects.
type GuildMembersResponse struct {
	Members []Member `json:"members"`
}

// It checks if a new member has joined the guild, and if so, it returns the new member's information
func onNewMemberInGuild(req static.AreaRequest) shared.AreaResponse {

	guildID := req.AuthStore["guild_id"]

	if (*req.Store)["req:guild:id"] != nil {
		guildID = (*req.Store)["req:guild:id"]
	}

	encode, _, err := req.Service.Endpoints["ListGuildMembersEndpoint"].CallEncode([]interface{}{
		guildID,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	memberList := GuildMembersResponse{}
	errr := json.Unmarshal(encode, &memberList.Members)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbMembers := len(memberList.Members)
	if (*req.Store)["ctx:members:number"] == nil {
		(*req.Store)["ctx:members:number"] = nbMembers
		(*req.Store)["ctx:members"] = memberList.Members
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:members:number"] == nbMembers {
		return shared.AreaResponse{Success: false}
	}

	if nbMembers < (*req.Store)["ctx:members:number"].(int) {
		(*req.Store)["ctx:members:number"] = nbMembers
		(*req.Store)["ctx:members"] = memberList.Members
		return shared.AreaResponse{Success: false}
	}

	if nbMembers <= 0 {
		return shared.AreaResponse{Error: errors.New("Invalid number of guilds")}
	}

	// Find the new guild in the list
	newMember := memberList.Members[0]
	for _, member := range memberList.Members {
		found := false
		for _, oldMember := range (*req.Store)["ctx:members"].([]Member) {
			if member.User.ID == oldMember.User.ID {
				found = true
				break
			}
		}
		if !found {
			newMember = member
			break
		}
	}

	if (*req.Store)["req:allow:bot"] != nil && (*req.Store)["req:allow:bot"].(string) == "false" && newMember.User.Bot {
		(*req.Store)["ctx:members:number"] = nbMembers
		(*req.Store)["ctx:members"] = memberList.Members
		return shared.AreaResponse{Success: false}
	}

	(*req.Store)["ctx:members:number"] = nbMembers
	(*req.Store)["ctx:members"] = memberList.Members
	req.Logger.WriteInfo("[Action] New member in guild (User: "+newMember.User.Username+"#"+newMember.User.Discriminator+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:guild:id":           guildID,
			"discord:user:id":            newMember.User.ID,
			"discord:user:username":      newMember.User.Username,
			"discord:user:discriminator": newMember.User.Discriminator,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForDiscordActionNewMemberInGuild() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_member_in_guild",
		Description: "When a new member is in a guild, (<1000 members only)",
		RequestStore: map[string]static.StoreElement{
			"req:allow:bot": {
				Priority:    1,
				Type:        "select",
				Description: "If true, the action will be triggered even if the member is a bot",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:guild:id": {
				Type:        "select_uri",
				Description: "The guild to check in the new channel (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
		},
		Method: onNewMemberInGuild,
		Components: []string{
			"discord:guild:id",
			"discord:user:id",
			"discord:user:username",
			"discord:user:discriminator",
		},
	}
}
