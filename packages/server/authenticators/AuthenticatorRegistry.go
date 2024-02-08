package authenticators

import (
	"area-server/authenticators/oauth2"
	"area-server/classes/static"
)

// It's a list of all the authenticators.
var List []static.OAuth2Authenticator = []static.OAuth2Authenticator{
	oauth2.DiscordAuthenticator(),
	oauth2.DropboxAuthenticator(),
	oauth2.FacebookAuthenticator(),
	oauth2.GithubAuthenticator(),
	oauth2.GoogleAuthenticator(),
	oauth2.MicrosoftAuthenticator(),
	oauth2.RedditAuthenticator(),
	oauth2.SpotifyAuthenticator(),
	oauth2.TwitchAuthenticator(),
}

// TODO: Add "state" parameter to all authenticators
// It returns the authenticator with the given name
func GetAuthenticator(name string) *static.OAuth2Authenticator {
	for _, authenticator := range List {
		if authenticator.Name == name {
			return &authenticator
		}
	}
	return nil
}
