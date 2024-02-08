package youtube

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
)

// It returns a map of strings to static.ServiceEndpoint structs
func YoutubeEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Validators
		"GetAllChannelsEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/channels",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Routes
		"GetAllPlaylistsEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/playlists",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"GetAllSubscriptionsEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/subscriptions",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllVideosEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/videos",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllRepliesEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/comments",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllCommentsEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/commentThreads",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetPlaylistItemsEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/playlistItems",
			Params:         GetPlaylistEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"AddSubscriptionEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/subscriptions",
			Params:         AddSubscriptionEndpointParams,
			ExpectedStatus: []int{200},
		},
		"AddCommentThreadEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/commentThreads",
			Params:         AddCommentEndpointParams,
			ExpectedStatus: []int{200},
		},
		"AddCommentEndpoint": {
			BaseURL:        "https://www.googleapis.com/youtube/v3/comments",
			Params:         AddCommentEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func BasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
		},
		QueryParams: params[1].(map[string]string),
	}
}

// It takes in a slice of interfaces and returns a pointer to a RequestParams struct
func GetPlaylistEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
		},
		QueryParams: map[string]string{
			"part":       "snippet",
			"maxResults": "50",
			"playlistId": params[1].(string),
		},
	}
}

// This function takes in two parameters, the first is a pointer to a models.Authorization struct and
// the second is a string. It returns a pointer to a utils.RequestParams struct.
func BasicPostEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		Body: params[1].(string),
	}
}

// It takes in two parameters, the first being a pointer to an `Authorization` struct and the second
// being a string. It returns a pointer to a `RequestParams` struct
func AddSubscriptionEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"part": "snippet",
		},
		Body: params[1].(string),
	}
}

// It takes two parameters, the first one is a pointer to a struct of type `models.Authorization` and
// the second one is a string. It returns a pointer to a struct of type `utils.RequestParams`
func AddCommentEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"part": "snippet",
		},
		Body: params[1].(string),
	}
}
