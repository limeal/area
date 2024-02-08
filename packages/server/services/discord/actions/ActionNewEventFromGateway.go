package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/discord/gateway"
	"encoding/json"
	"time"
)

// It waits for a new event to be received from the gateway, and if it is, it returns the event type
// and data
func onNewEventReceivedFromGateway(req static.AreaRequest) shared.AreaResponse {

	discordGateway := req.Service.Gateway.(*gateway.DiscordGateway)

	select {
	case event := <-discordGateway.EventChan:
		if (*req.Store)["req:gateway:event:type"] != nil && ((*req.Store)["req:gateway:event:type"].(string) != event.T) {
			return shared.AreaResponse{
				Success: false,
			}
		}

		data, err := json.Marshal(event.D)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		req.Logger.WriteInfo("[Action] New event received from gateway (Type: "+event.T+")", false)
		return shared.AreaResponse{
			Success: true,
			Data: map[string]interface{}{
				"discord:gateway:event:type": event.T,
				"discord:gateway:event:data": string(data),
			},
		}
	case <-time.After(100 * time.Millisecond):
		return shared.AreaResponse{
			Success: false,
		}
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForDiscordActionNewEventFromGateway() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_event_from_gateway",
		Description: "When a new event is received from the gateway",
		UseGateway:  true,
		RequestStore: map[string]static.StoreElement{
			"req:gateway:event:type": {
				Type:        "select",
				Description: "Only trigger when the event type is the one selected",
				Required:    false,
				Values: []string{
					"READY",
					"RESUMED",
					"CHANNEL_CREATE",
					"CHANNEL_UPDATE",
					"CHANNEL_DELETE",
					"CHANNEL_PINS_UPDATE",
					"GUILD_CREATE",
					"GUILD_UPDATE",
					"GUILD_DELETE",
					"GUILD_BAN_ADD",
					"GUILD_BAN_REMOVE",
					"GUILD_EMOJIS_UPDATE",
					"GUILD_INTEGRATIONS_UPDATE",
					"GUILD_MEMBER_ADD",
					"GUILD_MEMBER_REMOVE",
					"GUILD_MEMBER_UPDATE",
					"GUILD_MEMBERS_CHUNK",
					"GUILD_ROLE_CREATE",
					"GUILD_ROLE_UPDATE",
					"GUILD_ROLE_DELETE",
					"MESSAGE_CREATE",
					"MESSAGE_UPDATE",
					"MESSAGE_DELETE",
					"MESSAGE_DELETE_BULK",
					"MESSAGE_REACTION_ADD",
					"MESSAGE_REACTION_REMOVE",
					"MESSAGE_REACTION_REMOVE_ALL",
					"PRESENCE_UPDATE",
					"TYPING_START",
					"USER_UPDATE",
					"VOICE_STATE_UPDATE",
					"VOICE_SERVER_UPDATE",
					"WEBHOOKS_UPDATE",
				},
			},
		},
		Method: onNewEventReceivedFromGateway,
		Components: []string{
			"discord:gateway:event:type",
			"discord:gateway:event:data",
		},
	}
}
