package reddit

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"net/url"
)

// It returns a map of strings to static.ServiceEndpoint structs
func RedditEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Actions
		"GetUserEndpoint": {
			BaseURL:        "https://oauth.reddit.com/api/v1/me",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetSubredditEndpoint": {
			BaseURL:        "https://oauth.reddit.com/r/${subreddit}",
			Params:         GetSubredditEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserSubredditsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/subreddits/mine/subscriber",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"SearchSubredditsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/subreddits/search",
			Params:         SearchSubredditsEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetNewPostsOnSubRedditEndpoint": {
			BaseURL:        "https://oauth.reddit.com/r/${subreddit}/new",
			Params:         GetPostsOnSubRedditEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetHotPostsOnSubRedditEndpoint": {
			BaseURL:        "https://oauth.reddit.com/r/${subreddit}/hot",
			Params:         GetPostsOnSubRedditEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetTopPostsOnSubRedditEndpoint": {
			BaseURL:        "https://oauth.reddit.com/r/${subreddit}/top",
			Params:         GetPostsOnSubRedditEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserPostsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/user/${username}/submitted",
			Params:         GetUserDataEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserCommentsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/user/${username}/comments",
			Params:         GetUserDataEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserSavedPostsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/user/${username}/saved",
			Params:         GetUserDataEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserUpvotedPostsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/user/${username}/upvoted",
			Params:         GetUserDataEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserDownvotedPostsEndpoint": {
			BaseURL:        "https://oauth.reddit.com/user/${username}/downvoted",
			Params:         GetUserDataEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"ComposePrivateMessageEndpoint": {
			BaseURL:        "https://oauth.reddit.com/api/compose",
			Params:         GetComposePrivateMessageEndpointParams,
			ExpectedStatus: []int{200},
		},
		"SubmitNewLinkEndpoint": {
			BaseURL:        "https://oauth.reddit.com/api/submit",
			Params:         GetSubmitLinkEndpointParams,
			ExpectedStatus: []int{200},
		},
		"SubmitNewPostEndpoint": {
			BaseURL:        "https://oauth.reddit.com/api/submit",
			Params:         GetSubmitTextPostEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func GetBasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetSubredditEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
		},
		UrlParams: map[string]string{
			"subreddit": params[1].(string),
		},
		QueryParams: map[string]string{
			"limit": "1",
		},
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func SearchSubredditsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
		},
		QueryParams: map[string]string{
			"q":     params[1].(string),
			"sort":  "relevance",
			"limit": "1",
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetPostsOnSubRedditEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
		},
		UrlParams: map[string]string{
			"subreddit": params[1].(string),
		},
		QueryParams: map[string]string{
			"limit": "100",
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetUserDataEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
		},
		UrlParams: map[string]string{
			"username": params[1].(string),
		},
		QueryParams: params[2].(map[string]string),
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func GetComposePrivateMessageEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"to":      {params[1].(string)},
			"subject": {params[2].(string)},
			"text":    {params[3].(string)},
		}.Encode(),
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetSubmitLinkEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"sr":          {params[1].(string)},
			"title":       {params[2].(string)},
			"url":         {params[3].(string)},
			"kind":        {"link"},
			"sendreplies": {"true"},
		}.Encode(),
	}
}

// "This function returns a pointer to a RequestParams struct that contains the parameters for a POST
// request to the /api/submit endpoint with the given parameters."
//
// The first line is a comment that describes what the function does. The second line is the function
// signature. The third line is a comment that describes the parameters. The fourth line is a comment
// that describes the return value. The fifth line is the function body
func GetSubmitTextPostEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"User-Agent":    "Area:v1.0 (by /u/area)",
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		Body: url.Values{
			"sr":          {params[1].(string)},
			"title":       {params[2].(string)},
			"text":        {params[3].(string)},
			"kind":        {"self"},
			"sendreplies": {"true"},
		}.Encode(),
	}
}
