package oauth2

import (
	"area-server/classes/static"
	"area-server/utils"
	"net/url"
	"os"
)

// It takes in two parameters, the first being the authorization code and the second being the redirect
// URI, and returns a request params object that contains the method, headers, and body of the request
func getMicrosoftTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {params[0].(string)},
			"redirect_uri":  {params[1].(string)},
			"client_id":     {os.Getenv("MICROSOFT_CLIENT_ID")},
			"client_secret": {os.Getenv("MICROSOFT_SECRET_ID")},
		}.Encode(),
	}
}

// It takes a list of parameters, and returns a request params object
func getMicrosoftProfileRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + fields["access_token"].(string),
		},
	}
}

// It returns a URL that the user can visit to authorize the application to access their Microsoft
// account
func getMicrosoftAuthorizationURI() string {
	return "https://login.microsoftonline.com/organizations/oauth2/v2.0/authorize?client_id=" + os.Getenv("MICROSOFT_CLIENT_ID") + "&response_type=code&response_mode=query&scope=openid%20offline_access%20profile%20email%20User.Read"
}

// It returns a static.OAuth2Authenticator that uses the Microsoft Graph API to get the user's email
// address
func MicrosoftAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("MICROSOFT_CLIENT_ID")
	if !p {
		panic("MICROSOFT_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("MICROSOFT_SECRET_ID")
	if !pt {
		panic("MICROSOFT_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "microsoft",
		Enabled: true,
		More: static.More{
			Avatar: true,
			Color:  "#2f2f2f",
		},
		AuthorizationURI: getMicrosoftAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://login.microsoftonline.com/organizations/oauth2/v2.0/token",
				Params:         getMicrosoftTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Email: &utils.RequestDescriptor{
				BaseURL:        "https://graph.microsoft.com/v1.0/me",
				Params:         getMicrosoftProfileRequestParams,
				ExpectedStatus: []int{200},
				TransformResponse: func(response interface{}) (map[string]interface{}, error) {
					fields := response.(map[string]interface{})
					return map[string]interface{}{
						"email": fields["mail"].(string),
					}, nil
				},
			},
		},
	}
}
