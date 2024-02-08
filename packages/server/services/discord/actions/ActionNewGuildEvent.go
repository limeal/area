package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// `GuildEventsResponse` is a struct that contains a slice of structs that each contain a string, a
// string, a string, a string, a string, a string, a string, a string, a string, and a struct that
// contains a string, a string, a string, and a bool.
// @property {[]struct {
// 		ID          string `json:"id"`
// 		GuildID     string `json:"guild_id"`
// 		ChannelID   string `json:"channel_id"`
// 		CreatorID   string `json:"creator_id"`
// 		Name        string `json:"name"`
// 		Description string `json:"description"`
// 		StartTime   string `json:"start_time"`
// 		EndTime     string `json:"end_time"`
// 		Location    string `json:"location"`
// 		Creator     struct {
// 			ID            string `json:"id"`
// 			Username      string `json:"username"`
// 			Discriminator string `json:"discriminator"`
// 			Bot           bool   `json:"bot"`
// 		} `json:"creator"`
// 	}} Events - An array of events.
type GuildEventsResponse struct {
	Events []struct {
		ID          string `json:"id"`
		GuildID     string `json:"guild_id"`
		ChannelID   string `json:"channel_id"`
		CreatorID   string `json:"creator_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		Location    string `json:"location"`
		Creator     struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Discriminator string `json:"discriminator"`
			Bot           bool   `json:"bot"`
		} `json:"creator"`
	}
}

// It checks if the latest event in a guild is a new one, and if it is, it returns the event's
// information
func onNewGuildEvent(req static.AreaRequest) shared.AreaResponse {

	guildID := req.AuthStore["guild_id"]

	if (*req.Store)["req:guild:id"] != nil {
		guildID = (*req.Store)["req:guild:id"]
	}

	encode, _, err := req.Service.Endpoints["ListGuildEventsEndpoint"].CallEncode([]interface{}{
		guildID,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	eventList := GuildEventsResponse{}
	errr := json.Unmarshal(encode, &eventList.Events)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbEvents := len(eventList.Events)
	ok, errL := utils.IsLatestBasic(req.Store, nbEvents)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:allow:bot"] != nil && (*req.Store)["req:allow:bot"].(string) == "false" {
		if eventList.Events[0].Creator.Bot {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New event in guild (Name: "+eventList.Events[0].Name+") (Description: "+eventList.Events[0].Description+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"discord:guild:id":           guildID,
			"discord:event:id":           eventList.Events[0].ID,
			"discord:event:name":         eventList.Events[0].Name,
			"discord:event:description":  eventList.Events[0].Description,
			"discord:user:id":            eventList.Events[0].CreatorID,
			"discord:user:username":      eventList.Events[0].Creator.Username,
			"discord:user:discriminator": eventList.Events[0].Creator.Discriminator,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForDiscordActionNewEventInGuild() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_guild_event",
		Description: "When a new guild event is created",
		Method:      onNewGuildEvent,
		RequestStore: map[string]static.StoreElement{
			"req:guild:id": {
				Type:        "select_uri",
				Description: "The guild to check in the new event (default: guild selected in the auth)",
				Required:    false,
				Values:      []string{"/guilds?bot=true"},
			},
			"req:allow:bot": {
				Priority:    1,
				Type:        "select",
				Description: "Allow bot to trigger this action if it's the creator of the event",
				Required:    false,
				Values:      []string{"true", "false"},
			},
		},
		Components: []string{
			"discord:guild:id",
			"discord:event:id",
			"discord:event:name",
			"discord:event:description",
			"discord:user:id",
			"discord:user:username",
			"discord:user:discriminator",
		},
	}
}
