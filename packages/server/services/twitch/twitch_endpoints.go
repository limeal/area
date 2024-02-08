package twitch

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"os"
)

// It returns a map of strings to static.ServiceEndpoint structs
func TwitchEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Routes
		"GetGameEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/games",
			Params:         GetGameEndpointParams,
			ExpectedStatus: []int{200},
		},
		// SERVICE SPECIFIC
		"GetUserByIdEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/users",
			Params:         GetUserByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserByLoginEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/users",
			Params:         GetUserByLoginEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserFollowsEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/users/follows",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetSoundTrackCurrentTrackEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/soundtrack/current_track",
			Params:         GetSoundTrackCurrentTrackEndpointParams,
			ExpectedStatus: []int{200, 404},
		},
		"GetClipsEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/clips",
			Params:         GetClipsEndpointParams,
			ExpectedStatus: []int{200},
		},
		"NewStreamStartedByUserEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/streams",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateClipEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/clips",
			Params:         CreateClipEndpointParams,
			ExpectedStatus: []int{202},
		},
		"SendWhisperEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/whispers",
			Params:         SendWhisperEndpointParams,
			ExpectedStatus: []int{204},
		},
		"CreatePollEndpoint": {
			BaseURL:        "https://api.twitch.tv/helix/polls",
			Params:         CreatePollEndpointParams,
			ExpectedStatus: []int{201},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func GetUserByIDEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
		},
		QueryParams: map[string]string{
			"id": params[1].(string),
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetUserByLoginEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
		},
		QueryParams: map[string]string{
			"login": params[1].(string),
		},
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func GetGameEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Client-ID":     os.Getenv("TWITCH_CLIENT_ID"),
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
		},
		QueryParams: map[string]string{
			"name": params[1].(string),
		},
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func GetClipsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
		},
		QueryParams: params[1].(map[string]string),
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func GetBasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
		},
		QueryParams: params[1].(map[string]string),
	}
}

// It creates a request params object for the CreateClip endpoint
func CreateClipEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"broadcaster_id": params[1].(string),
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetSoundTrackCurrentTrackEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
		},
		QueryParams: map[string]string{
			"broadcaster_id": params[1].(string),
		},
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func SendWhisperEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"from_user_id": params[1].(string),
			"to_user_id":   params[2].(string),
		},
		Body: "{\"message\":\"" + params[3].(string) + "\"}",
	}
}

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func CreatePollEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Client-Id":     os.Getenv("TWITCH_CLIENT_ID"),
			"Content-Type":  "application/json",
		},
		Body: params[1].(string),
	}
}
