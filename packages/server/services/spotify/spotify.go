package spotify

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/spotify/actions"
	"area-server/services/spotify/reactions"
)

// It returns a static.Service object that describes the Spotify service
func Descriptor() static.Service {
	return static.Service{
		Name:          "spotify",
		Description:   "Spotify is a music, podcast, and video streaming service that gives you access to millions of songs and other content from artists all over the world.",
		Authenticator: authenticators.GetAuthenticator("spotify"),
		RateLimit:     3,
		Validators:    SpotifyValidators(),
		Endpoints:     SpotifyEndpoints(),
		Routes:        SpotifyRoutes(),
		Actions: []static.ServiceArea{
			actions.DescriptorForSpotifyActionNewSavedTrack(),
			actions.DescriptorForSpotifyActionNewTrackAddedToPlaylist(),
			actions.DescriptorForSpotifyActionNewSavedAlbum(),
			actions.DescriptorForSpotifyActionNewSavedShow(),
			actions.DescriptorForSpotifyActionNewSavedEpisode(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForSpotifyReactionAddItemToPlaylist(),
			reactions.DescriptorForSpotifyReactionAddItemToPlaylistFromSearch(),
			reactions.DescriptorForSpotifyReactionPausePlayback(),
			reactions.DescriptorForSpotifyReactionSkipToNextPlayback(),
			reactions.DescriptorForSpotifyReactionSkipToPreviousPlayback(),
			reactions.DescriptorForSpotifyReactionResumePlayback(),
			reactions.DescriptorForSpotifyReactionAddItemToPlaybackQueue(),
		},
	}
}
