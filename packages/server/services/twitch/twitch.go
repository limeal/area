package twitch

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/twitch/actions"
	"area-server/services/twitch/reactions"
)

// It returns a `static.Service` object that describes the Twitch service
func Descriptor() static.Service {
	return static.Service{
		Name:          "twitch",
		Description:   "Twitch is a live streaming video platform owned by Twitch Interactive, a subsidiary of Amazon.",
		Authenticator: authenticators.GetAuthenticator("twitch"),
		RateLimit:     4,
		Endpoints:     TwitchEndpoints(),
		Validators:    TwitchValidators(),
		Actions: []static.ServiceArea{
			actions.DescriptorForTwitchActionFollowNewStreamer(),
			actions.DescriptorForTwitchActionNewStreamStarted(),
			actions.DescriptorForTwitchActionNewClipCaptured(),
			actions.DescriptorForTwitchActionCurrentTrackChange(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForTwitchReactionCreateClip(),
			reactions.DescriptorForTwitchReactionSendWhisper(),
			reactions.DescriptorForTwitchReactionCreatePoll(),
		},
	}
}
