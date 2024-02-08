package oauth2

import (
	"area-server/classes/static"
	"area-server/utils"
	"os"
	"strings"
)

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func getFacebookTokenRequestParams(params []interface{}) *utils.RequestParams {
	redirectURI := params[1].(string)

	// If redirect uri doesn't end with a slash, add one
	if redirectURI[len(redirectURI)-1] != '/' {
		redirectURI += "/"
	}

	return &utils.RequestParams{
		Method: "GET",
		QueryParams: map[string]string{
			"client_id":     os.Getenv("FACEBOOK_CLIENT_ID"),
			"client_secret": os.Getenv("FACEBOOK_SECRET_ID"),
			"code":          params[0].(string),
			"redirect_uri":  redirectURI,
		},
	}
}

// It takes a slice of interfaces, and returns a pointer to a RequestParams struct
func getFacebookEmailRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		QueryParams: map[string]string{
			"access_token": fields["access_token"].(string),
			"fields":       "email",
		},
	}
}

// It returns the Facebook authorization URI
func getFacebookAuthorizationURI() string {
	return "https://www.facebook.com/v16.0/dialog/oauth?client_id=" + os.Getenv("FACEBOOK_CLIENT_ID") + "&response_type=code&scope=" + strings.Join([]string{
		"email",
		// User
		"user_birthday",
		"user_friends",
		"user_gender",
		"user_hometown",
		"user_likes",
		"user_location",
		"user_photos",
	}, "%20")
}

// It returns a static.OAuth2Authenticator struct that contains the information needed to authenticate
// with Facebook
func FacebookAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("FACEBOOK_CLIENT_ID")
	if !p {
		panic("FACEBOOK_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("FACEBOOK_SECRET_ID")
	if !pt {
		panic("FACEBOOK_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "facebook",
		Enabled: true,
		More: static.More{
			Avatar: true,
			Color:  "#3b5998",
		},
		AuthorizationURI: getFacebookAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://graph.facebook.com/v16.0/oauth/access_token",
				Params:         getFacebookTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Email: &utils.RequestDescriptor{
				BaseURL:        "https://graph.facebook.com/v16.0/me",
				Params:         getFacebookEmailRequestParams,
				ExpectedStatus: []int{200},
			},
		},
	}
}
