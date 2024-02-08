package oauth2

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"net/url"
	"os"
	"strings"
)

// It creates a request to the Discord API to get a token
func getDiscordTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
			"client_secret": {os.Getenv("DISCORD_SECRET_ID")},
			"grant_type":    {"authorization_code"},
			"code":          {params[0].(string)},
			"state":         {os.Getenv("AREA_STATE")},
			"redirect_uri":  {params[1].(string)},
		}.Encode(),
	}
}

// It returns a pointer to a `RequestParams` struct that contains the parameters for a POST request to
// the Discord API to get a new access token using a refresh token
func getDiscordRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
			"client_secret": {os.Getenv("DISCORD_SECRET_ID")},
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
		}.Encode(),
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func getDiscordRevokeTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
			"client_secret": {os.Getenv("DISCORD_SECRET_ID")},
			"token":         {params[0].(models.Authorization).AccessToken},
		}.Encode(),
	}
}

// It takes a list of parameters, and returns a request params object
func getDiscordProfileRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + fields["access_token"].(string),
		},
	}
}

// It returns a URL that you can use to authorize your bot to join your server
func getDiscordAuthorizationURI() string {
	return "https://discordapp.com/api/oauth2/authorize?client_id=" + os.Getenv("DISCORD_CLIENT_ID") + "&response_type=code&permissions=8&scope=" + strings.Join([]string{
		// OAuth2
		"bot",
		"identify",
		"email",
		"connections",
		// Guilds
		"guilds",
		"guilds.join",
		"guilds.members.read",
		// Messages
		"messages.read",
		// Other
		"role_connections.write",
	}, "%20") + "&state=" + os.Getenv("AREA_STATE")
}

// It returns a static.OAuth2Authenticator struct with the Name, Enabled, More, AuthorizationURI,
// AuthEndpoints, and OtherParams fields set
func DiscordAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("DISCORD_CLIENT_ID")
	if !p {
		panic("DISCORD_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("DISCORD_SECRET_ID")
	if !pt {
		panic("DISCORD_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "discord",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#7289DA",
		},
		AuthorizationURI: getDiscordAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://discordapp.com/api/oauth2/token",
				Params:         getDiscordTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://discordapp.com/api/oauth2/token",
				Params:         getDiscordRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RevokeToken: &utils.RequestDescriptor{
				BaseURL:        "https://discordapp.com/api/oauth2/token/revoke",
				Params:         getDiscordRevokeTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Profile: &utils.RequestDescriptor{
				BaseURL:        "https://discordapp.com/api/users/@me",
				Params:         getDiscordProfileRequestParams,
				ExpectedStatus: []int{200},
			},
		},
		OtherParams: func(m map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"guild_id": m["guild"].(map[string]interface{})["id"],
				"id":       m["id"],
			}
		},
	}
}
