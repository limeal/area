package oauth2

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"encoding/base64"
	"net/url"
	"os"
	"strings"
)

// It returns a pointer to a `RequestParams` struct with the method, headers, and body for the request
// to get a Reddit token
func getRedditTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("REDDIT_CLIENT_ID")+":"+os.Getenv("REDDIT_SECRET_ID"))),
		},
		Body: url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {params[0].(string)},
			"redirect_uri": {params[1].(string)},
			"state":        {"123456789"},
		}.Encode(),
	}
}

// It returns a pointer to a `RequestParams` struct that contains the method, headers, and body of the
// request
func getRedditRefreshTokenRequestParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("REDDIT_CLIENT_ID")+":"+os.Getenv("REDDIT_SECRET_ID"))),
		},
		Body: url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {params[0].(string)},
		}.Encode(),
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func getRedditRevokeTokenParams(params []interface{}) *utils.RequestParams {
	resp, err := utils.MakeRequest("https://www.reddit.com/api/v1/revoke_token", &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("REDDIT_CLIENT_ID")+":"+os.Getenv("REDDIT_SECRET_ID"))),
		},
		Body: url.Values{
			"token":           {params[0].(models.Authorization).RefreshToken},
			"token_type_hint": {"refresh_token"},
		}.Encode(),
	})

	if err != nil || resp.StatusCode != 200 {
		return nil
	}

	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("REDDIT_CLIENT_ID")+":"+os.Getenv("REDDIT_SECRET_ID"))),
		},
		Body: url.Values{
			"token":           {params[0].(models.Authorization).AccessToken},
			"token_type_hint": {"access_token"},
		}.Encode(),
	}
}

// It takes a slice of interfaces, and returns a pointer to a RequestParams struct
func getRedditProfileParams(params []interface{}) *utils.RequestParams {
	fields := params[0].(map[string]interface{})
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + fields["access_token"].(string),
			"User-Agent":    "Area:v1.0 (by /u/area)",
			"Accept":        "application/json",
		},
	}
}

// It returns a URL that you can use to get a Reddit access token
func getRedditAuthorizationURI() string {
	return "https://www.reddit.com/api/v1/authorize?client_id=" + os.Getenv("REDDIT_CLIENT_ID") + "&response_type=code&duration=permanent&scope=" + strings.Join([]string{
		"identity",
		"mysubreddits",
		"read",
		"save",
		"submit",
		"subscribe",
		"vote",
		"wikiread",
		"edit",
		"flair",
		"history",
		"modconfig",
		"modflair",
		"modlog",
		"modposts",
		"modwiki",
		"privatemessages",
	}, "%20") + "&state=123456789"
}

// It returns a static.OAuth2Authenticator struct that contains all the information needed to
// authenticate a user with Reddit
func RedditAuthenticator() static.OAuth2Authenticator {

	_, p := os.LookupEnv("REDDIT_CLIENT_ID")
	if !p {
		panic("REDDIT_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("REDDIT_SECRET_ID")
	if !pt {
		panic("REDDIT_SECRET_ID is not set")
	}

	return static.OAuth2Authenticator{
		Name:    "reddit",
		Enabled: false,
		More: static.More{
			Avatar: true,
			Color:  "#FF4500",
		},
		AuthorizationURI: getRedditAuthorizationURI(),
		AuthEndpoints: static.AuthEndpoints{
			AccessToken: utils.RequestDescriptor{
				BaseURL:        "https://www.reddit.com/api/v1/access_token",
				Params:         getRedditTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RefreshToken: &utils.RequestDescriptor{
				BaseURL:        "https://www.reddit.com/api/v1/access_token",
				Params:         getRedditRefreshTokenRequestParams,
				ExpectedStatus: []int{200},
			},
			RevokeToken: &utils.RequestDescriptor{
				BaseURL:        "https://www.reddit.com/api/v1/revoke_token",
				Params:         getRedditRevokeTokenParams,
				ExpectedStatus: []int{200},
			},
			Profile: &utils.RequestDescriptor{
				BaseURL:        "https://oauth.reddit.com/api/v1/me",
				Params:         getRedditProfileParams,
				ExpectedStatus: []int{200},
			},
		},
		OtherParams: func(m map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"name": m["name"],
			}
		},
	}
}
