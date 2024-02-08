package oauth2

import (
	"area-server/classes/static"
	"area-server/utils"
	"encoding/base64"
	"net/url"
	"os"
	"strings"
)

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func getSpotifyTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID")+":"+os.Getenv("SPOTIFY_SECRET_ID"))),
		},
		Body: url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {params[0].(string)},
			"redirect_uri": {params[1].(string)},
		}.Encode(),
	}
}

// It takes a slice of interfaces as a parameter, and returns a pointer to a RequestParams struct
func getSpotifyProfileRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + fields["access_token"].(string),
		},
	}
}

// It returns a pointer to a `utils.RequestParams` struct that contains the parameters for a POST
// request to the Spotify API to get a new access token using a refresh token
func getSpotifyRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID")+":"+os.Getenv("SPOTIFY_SECRET_ID"))),
		},
		Body: url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
		}.Encode(),
	}
}

// Fields is used to change the picked fields if they do not follow the default pattern

func getSpotifyAuthorizationURI() string {
	return "https://accounts.spotify.com/authorize?client_id=" + os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&scope=" + strings.Join([]string{
		// User
		"user-read-private",
		"user-read-email",
		// Playlist
		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-public",
		"playlist-modify-private",
		// History
		"user-top-read",
		"user-read-recently-played",
		"user-read-playback-position",
		// Follow
		"user-follow-modify",
		"user-follow-read",
		// Library
		"user-library-modify",
		"user-library-read",
		// Spotify Connect
		"user-read-playback-state",
		"user-modify-playback-state",
		"user-read-currently-playing",
	}, "%20") + "&show_dialog=true"
}

// It returns a static.OAuth2Authenticator struct that contains all the information needed to
// authenticate a user with Spotify
func SpotifyAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("SPOTIFY_CLIENT_ID")
	if !p {
		panic("SPOTIFY_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("SPOTIFY_SECRET_ID")
	if !pt {
		panic("SPOTIFY_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "spotify",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#1DB954",
		},
		AuthorizationURI: getSpotifyAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://accounts.spotify.com/api/token",
				Params:         getSpotifyTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://accounts.spotify.com/api/token",
				Params:         getSpotifyRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Profile: &utils.RequestDescriptor{
				BaseURL:        "https://api.spotify.com/v1/me",
				Params:         getSpotifyProfileRequestParams,
				ExpectedStatus: []int{200},
			},
		},
	}
}
