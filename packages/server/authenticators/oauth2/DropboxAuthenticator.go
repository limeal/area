package oauth2

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"net/url"
	"os"
)

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func getDropboxTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("DROPBOX_CLIENT_ID")},
			"client_secret": {os.Getenv("DROPBOX_SECRET_ID")},
			"grant_type":    {"authorization_code"},
			"code":          {params[0].(string)},
			"redirect_uri":  {params[1].(string)},
		}.Encode(),
	}
}

// It takes a single parameter, a string, and returns a pointer to a `utils.RequestParams` struct with
// the appropriate values set
func getDropboxRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"client_id":     {os.Getenv("DROPBOX_CLIENT_ID")},
			"client_secret": {os.Getenv("DROPBOX_SECRET_ID")},
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
		}.Encode(),
	}
}

// It returns a pointer to a `utils.RequestParams` struct with the `Method` field set to `"POST"` and
// the `Headers` field set to a map with a single key-value pair, where the key is `"Authorization"`
// and the value is `"Bearer "` concatenated with the `AccessToken` field of the `models.Authorization`
// struct passed in as the first parameter
func getDropboxRevokeTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(models.Authorization).AccessToken,
		},
	}
}

// It returns the authorization URI for Dropbox
func getDropboxAuthorizationURI() string {
	return "https://www.dropbox.com/oauth2/authorize?client_id=" + os.Getenv("DROPBOX_CLIENT_ID") + "&response_type=code&token_access_type=offline&state=" + os.Getenv("AREA_STATE")
}

// It returns a static.OAuth2Authenticator struct with the name "dropbox", the color #007EE5, and the
// authorization URI, access token, refresh token, and revoke token endpoints
func DropboxAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("DROPBOX_CLIENT_ID")
	if !p {
		panic("DROPBOX_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("DROPBOX_SECRET_ID")
	if !pt {
		panic("DROPBOX_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "dropbox",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#007EE5",
		},
		AuthorizationURI: getDropboxAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://api.dropboxapi.com/oauth2/token",
				Params:         getDropboxTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://api.dropboxapi.com/oauth2/token",
				Params:         getDropboxRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RevokeToken: &utils.RequestDescriptor{
				BaseURL:        "https://api.dropboxapi.com/2/auth/token/revoke",
				Params:         getDropboxRevokeTokenRequestParams,
				ExpectedStatus: []int{200},
			},
		},
	}
}
