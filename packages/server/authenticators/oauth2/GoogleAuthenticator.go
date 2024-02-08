package oauth2

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"net/url"
	"os"
	"strings"
)

// It takes in a slice of interfaces and returns a pointer to a RequestParams struct
func getGoogleAccessTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {params[0].(string)},
			"redirect_uri":  {params[1].(string)},
			"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
			"client_secret": {os.Getenv("GOOGLE_SECRET_ID")},
		}.Encode(),
	}
}

// It takes a refresh token as a parameter and returns a request to Google's API to get a new access
// token
func getGoogleRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
			"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
			"client_secret": {os.Getenv("GOOGLE_SECRET_ID")},
		}.Encode(),
	}
}

// It takes a slice of interfaces as an argument, and returns a pointer to a RequestParams struct
func getGoogleProfileRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + fields["access_token"].(string),
		},
	}
}

// It returns a pointer to a `utils.RequestParams` struct with the method, headers, and body set to the
// appropriate values for the Google revoke token request
func getGoogleRevokeTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"token": {params[0].(models.Authorization).AccessToken},
		}.Encode(),
	}
}

// It returns the URL that the user should be redirected to in order to authenticate with Google
func getGoogleAuthorizationURI() string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + os.Getenv("GOOGLE_CLIENT_ID") + "&response_type=code&scope=" + strings.Join([]string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/gmail.settings.basic",
		"https://mail.google.com/", // Gmail
		// Youtube
		"https://www.googleapis.com/auth/youtube",
		"https://www.googleapis.com/auth/youtube.force-ssl",
		"https://www.googleapis.com/auth/youtube.readonly",
		"https://www.googleapis.com/auth/youtube.upload",
		"https://www.googleapis.com/auth/youtubepartner",
		"https://www.googleapis.com/auth/youtubepartner-channel-audit",
	}, "%20") + "&access_type=offline&prompt=consent"
}

// It returns a static.OAuth2Authenticator struct with the name "google", the color #4285F4, and the
// endpoints for the authorization, access token, refresh token, revoke token, and email requests
func GoogleAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !p {
		panic("GOOGLE_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("GOOGLE_SECRET_ID")
	if !pt {
		panic("GOOGLE_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "google",
		Enabled: true,
		More: static.More{
			Avatar: true,
			Color:  "#4285F4",
		},
		AuthorizationURI: getGoogleAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://www.googleapis.com/oauth2/v4/token",
				Params:         getGoogleAccessTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://www.googleapis.com/oauth2/v4/token",
				Params:         getGoogleRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RevokeToken: &utils.RequestDescriptor{
				BaseURL:        "https://oauth2.googleapis.com/revoke",
				Params:         getGoogleRevokeTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Email: &utils.RequestDescriptor{
				BaseURL:        "https://www.googleapis.com/oauth2/v3/userinfo",
				Params:         getGoogleProfileRequestParams,
				ExpectedStatus: []int{200},
			},
		},
	}
}
