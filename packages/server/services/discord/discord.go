package discord

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/discord/actions"
	"area-server/services/discord/gateway"
	"area-server/services/discord/reactions"
	"os"
)

// It returns a static.Service object that describes the service
func Descriptor() static.Service {

	_, ptt := os.LookupEnv("DISCORD_BOT_TOKEN")
	if !ptt {
		panic("DISCORD_BOT_TOKEN is not set")
	}

	return static.Service{
		Name:          "discord",
		Description:   "Discord is a proprietary freeware VoIP application and digital distribution platform designed for video gaming communities, that specializes in text, image, video and audio communication between users in a chat channel.",
		Authenticator: authenticators.GetAuthenticator("discord"),
		RateLimit:     3,
		Endpoints:     DiscordEndpoints(),
		Validators:    DiscordValidators(),
		Routes:        DiscordRoutes(),
		Gateway: &gateway.DiscordGateway{
			EventChan: make(chan gateway.Event),
			Interrupt: make(chan bool, 1),
		},
		Actions: []static.ServiceArea{
			actions.DescriptorForDiscordActionNewMessageInChannel(),
			actions.DescriptorForDiscordActionNewGuildJoined(),
			actions.DescriptorForDiscordActionNewMemberInGuild(),
			actions.DescriptorForDiscordActionNewEventInGuild(),
			actions.DescriptorForDiscordActionNewPinnedMessage(),
			actions.DescriptorForDiscordActionNewReactionAddedByUser(),
			actions.DescriptorForDiscordActionNewChannelInGuild(),
			// Gateway
			actions.DescriptorForDiscordActionNewEventFromGateway(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForDiscordReactionPostMessage(),
			reactions.DescriptorForDiscordReactionSendPrivateMessage(),
			reactions.DescriptorForDiscordReactionAddReactionToMessage(),
			reactions.DescriptorForDiscordReactionBanUser(),
			reactions.DescriptorForDiscordReactionKickUser(),
			reactions.DescriptorForDiscordReactionPostEmbedMessage(),
		},
	}
}
