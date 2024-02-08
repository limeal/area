package oauth2

import (
	"area-server/classes/static"
	"area-server/utils"
	"os"
	"strings"
)

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func getGithubTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		QueryParams: map[string]string{
			"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
			"client_secret": os.Getenv("GITHUB_SECRET_ID"),
			"code":          params[0].(string),
			"state":         os.Getenv("AREA_STATE"),
			"redirect_uri":  params[1].(string),
		},
	}
}

// It takes a list of parameters, and returns a request params object
func getGithubProfileRequestParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Accept":               "application/vnd.github.v3+json",
			"Authorization":        "Bearer " + fields["access_token"].(string),
			"X-GitHub-Api-Version": "2022-11-28",
		},
	}
}

// It returns the URI to redirect the user to in order to get a GitHub authorization code
func getGithubAuthorizationURI() string {
	return "https://github.com/login/oauth/authorize?client_id=" + os.Getenv("GITHUB_CLIENT_ID") + "&scope=" + strings.Join([]string{
		"user",
		"repo",
		"notifications",
		"gist",
		"read:user",
		"user:email",
		"user:follow",
	}, "%20") + "&state=" + os.Getenv("AREA_STATE")
}

// It returns a static.OAuth2Authenticator that uses the Github API to authenticate users
func GithubAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("GITHUB_CLIENT_ID")
	if !p {
		panic("GITHUB_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("GITHUB_SECRET_ID")
	if !pt {
		panic("GITHUB_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "github",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#24292e",
		},
		AuthorizationURI: getGithubAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://github.com/login/oauth/access_token",
				Params:         getGithubTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			Email: &utils.RequestDescriptor{
				BaseURL:        "https://api.github.com/user/emails",
				Params:         getGithubProfileRequestParams,
				ExpectedStatus: []int{200},
				TransformResponse: func(response any) (map[string]interface{}, error) {
					emails := response.([]interface{})
					return map[string]interface{}{
						"email": emails[0].(map[string]interface{})["email"],
					}, nil
				},
			},
			Profile: &utils.RequestDescriptor{
				BaseURL:        "https://api.github.com/user",
				Params:         getGithubProfileRequestParams,
				ExpectedStatus: []int{200},
				TransformResponse: func(response any) (map[string]interface{}, error) {
					fields := response.(map[string]interface{})
					return map[string]interface{}{
						"login": fields["login"],
					}, nil
				},
			},
		},
		OtherParams: func(m map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"login": m["login"],
			}
		},
	}
}
