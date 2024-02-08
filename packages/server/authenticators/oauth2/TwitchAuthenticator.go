package oauth2

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"net/url"
	"os"
	"strings"
)

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func getTwitchTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("TWITCH_CLIENT_ID")},
			"client_secret": {os.Getenv("TWITCH_SECRET_ID")},
			"code":          {params[0].(string)},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {params[1].(string)},
		}.Encode(),
	}
}

// It takes a single parameter, a string, and returns a pointer to a RequestParams struct with the
// appropriate values for a Twitch API refresh token request
func getTwitchRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("TWITCH_CLIENT_ID")},
			"client_secret": {os.Getenv("TWITCH_SECRET_ID")},
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
		}.Encode(),
	}
}

// It takes a slice of interfaces as a parameter, and returns a pointer to a `RequestParams` struct
func getTwitchValidateTokenRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + fields["access_token"].(string),
		},
	}
}

// It returns a pointer to a `RequestParams` struct that contains the parameters for the request to the
// Twitch API to revoke the access token
func getTwitchRevokeTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Accept":       "*/*",
		},
		Body: url.Values{
			"client_id": {os.Getenv("TWITCH_CLIENT_ID")},
			"token":     {params[0].(models.Authorization).AccessToken},
		}.Encode(),
	}
}

// It returns a URL that you can use to get a user's authorization code
func getTwitchAuthorizationURI() string {
	return "https://id.twitch.tv/oauth2/authorize?client_id=" + os.Getenv("TWITCH_CLIENT_ID") + "&response_type=code&scope=" + strings.Join([]string{
		// Channel
		"channel:manage:broadcast",
		"channel:manage:moderators",
		"channel:manage:polls",
		"channel:manage:predictions",
		"channel:manage:raids",
		"channel:manage:redemptions",
		"channel:manage:schedule",
		"channel:manage:videos",
		"channel:read:editors",
		"channel:read:goals",
		"channel:read:hype_train",
		"channel:read:polls",
		"channel:read:predictions",
		"channel:read:redemptions",
		"channel:read:stream_key",
		"channel:read:subscriptions",
		"channel:manage:extensions",
		"channel:read:subscriptions",
		"channel:read:vips",
		"channel:manage:vips",
		// Clips
		"clips:edit",
		// Moderation
		"moderation:read",
		"moderator:manage:announcements",
		"moderator:manage:automod",
		"moderator:manage:banned_users",
		"moderator:manage:blocked_terms",
		"moderator:read:blocked_terms",
		"moderator:manage:chat_messages",
		// Chat
		"chat:edit",
		"chat:read",
		// User
		"user:edit",
		"user:edit:broadcast",
		"user:manage:blocked_users",
		"user:read:blocked_users",
		"user:read:broadcast",
		"user:read:email",
		"user:read:follows",
		"user:read:subscriptions",
		"user:manage:whispers",
	}, "%20")
}

// It returns a static.OAuth2Authenticator that uses the Twitch OAuth2 API to authenticate users
func TwitchAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("TWITCH_CLIENT_ID")
	if !p {
		panic("TWITCH_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("TWITCH_SECRET_ID")
	if !pt {
		panic("TWITCH_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "twitch",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#6441A4",
		},
		AuthorizationURI: getTwitchAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://id.twitch.tv/oauth2/token",
				Params:         getTwitchTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://id.twitch.tv/oauth2/token",
				Params:         getTwitchRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RevokeToken: &utils.RequestDescriptor{
				BaseURL:        "https://id.twitch.tv/oauth2/revoke",
				Params:         getTwitchRevokeTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			ValidateToken: &utils.RequestDescriptor{
				BaseURL:        "https://id.twitch.tv/oauth2/validate",
				Params:         getTwitchValidateTokenRequestParams,
				ExpectedStatus: []int{200},
				TransformResponse: func(response any) (map[string]interface{}, error) {
					fields := response.(map[string]interface{})
					return map[string]interface{}{
						"user_id": fields["user_id"],
						"login":   fields["login"],
					}, nil
				},
			},
		},
		OtherParams: func(m map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"user_id": m["user_id"],
				"login":   m["login"],
			}
		},
	}
}
