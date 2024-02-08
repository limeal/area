package youtube

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/youtube/actions"
	"area-server/services/youtube/reactions"
)

// It returns a static.Service object that describes the service
func Descriptor() static.Service {
	return static.Service{
		Name:          "youtube",
		Description:   "YouTube is an American online video-sharing platform headquartered in San Bruno, California. Three former PayPal employees—Chad Hurley, Steve Chen, and Jawed Karim—created the service in February 2005. Google bought the site in November 2006 for US$1.65 billion; YouTube now operates as one of Google's subsidiaries.",
		Authenticator: authenticators.GetAuthenticator("google"),
		RateLimit:     2,
		More: &static.More{
			Avatar: true,
			Color:  "#FF0000",
		},
		Endpoints:  YoutubeEndpoints(),
		Validators: YoutubeValidators(),
		Routes:     YoutubeRoutes(),
		Actions: []static.ServiceArea{
			actions.DescriptorForYoutubeActionNewCommentByYou(),
			actions.DescriptorForYoutubeActionNewLikedVideo(),
			actions.DescriptorForYoutubeActionNewSubscriber(),
			actions.DescriptorForYoutubeActionNewSubscription(),
			actions.DescriptorForYoutubeActionNewVideoAddedToPlaylist(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForYoutubeReactionAddComment(),
			reactions.DescriptorForYoutubeReactionAddSubscription(),
		},
	}
}
