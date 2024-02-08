package spotify

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
)

// It returns a map of endpoint names to endpoint definitions
func SpotifyEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Routes
		"GetUserProfileEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me",
			Params:         GetUserProfileEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserPlaylistsEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/playlists",
			Params:         GetUserPlaylistsEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"FindTrackByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/tracks/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindShowByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/shows/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindEpisodeByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/episodes/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindArtistByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/artists/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindAlbumByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/albums/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindPlaylistByIDEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/playlists/${id}",
			Params:         FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindUserSavedTracksEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/tracks",
			Params:         FindUserSavedEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindUserSavedShowsEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/shows",
			Params:         FindUserSavedEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindUserSavedEpisodesEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/episodes",
			Params:         FindUserSavedEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindUserSavedAlbumsEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/albums",
			Params:         FindUserSavedEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserAvailableDevicesEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/player/devices",
			Params:         GetUserAvailableDevicesEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"AddItemToPlaybackQueue": {
			BaseURL:        "https://api.spotify.com/v1/me/player/queue",
			Params:         AddItemToPlaybackQueueParams,
			ExpectedStatus: []int{204},
		},
		"AddItemToPlaylistEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/playlists/${playlist_id}/tracks",
			Params:         AddItemToPlaylistEndpointParams,
			ExpectedStatus: []int{201},
		},
		"ResumePlaybackEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/player/play",
			Params:         ResumePlaybackEndpointParams,
			ExpectedStatus: []int{204},
		},
		"PausePlaybackEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/player/pause",
			Params:         PauseEndpointParams,
			ExpectedStatus: []int{204},
		},
		"SkipToNextPlaybackEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/player/next",
			Params:         DevideIDEndpointParams,
			ExpectedStatus: []int{204},
		},
		"SkipToPreviousPlaybackEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/me/player/previous",
			Params:         DevideIDEndpointParams,
			ExpectedStatus: []int{204},
		},
		// Other
		"SearchItemsEndpoint": {
			BaseURL:        "https://api.spotify.com/v1/search",
			Params:         SearchItemsEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func GetUserProfileEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
	}
}

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func GetUserPlaylistsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
	}
}

// --- FIND BY ID (2 params - [0] = auth, [1] = id)

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func FindByIDEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
	}
}

// --- AVAILABLE DEVICES (1 params - [0] = auth)

// It returns a `RequestParams` object with the `Method` set to `GET`, and the `Headers` set to a map
// with the `Authorization` and `Content-Type` headers
func GetUserAvailableDevicesEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
	}
}

// --- SAVED (1 param - [0] = auth)

// It takes a slice of interface{} and returns a pointer to a RequestParams struct
func FindUserSavedEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"limit": "1",
		},
	}
}

// --- PLAYBACK QUEUE (2 params - [0] = auth, [1] = query{uri, device_id})

// It takes in two parameters, the first being a pointer to a `models.Authorization` struct and the
// second being a map of strings to strings. It returns a pointer to a `utils.RequestParams` struct
func AddItemToPlaybackQueueParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: params[1].(map[string]string),
	}
}

// --- PLAYLIST TRACKS (3 params - [0] = auth, [1] = playlist_id, [2] = body{uris, position})

// It takes in an array of interfaces, and returns a pointer to a RequestParams struct
func AddItemToPlaylistEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"playlist_id": params[1].(string),
		},
		Body: params[2].(string),
	}
}

// --- RESUME PLAYBACK (2 params - [0] = auth, [1] = query{device_id}, [2] = body{uris, context_uri})

// It takes in three parameters, the first of which is a pointer to a `models.Authorization` struct,
// the second of which is a string, and the third of which is a string. It returns a pointer to a
// `utils.RequestParams` struct
func ResumePlaybackEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"device_id": params[1].(string),
		},
		Body: params[2].(string),
	}
}

// ---- PAUSE
// It takes two parameters, the first is a pointer to a `models.Authorization` struct and the second is
// a string. It returns a pointer to a `utils.RequestParams` struct
func PauseEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"device_id": params[1].(string),
		},
	}
}

// ---  SKIP TO NEXT / PREVIOUS (2 params - [0] = auth, [1] = query{device_id})

// It takes in a slice of interfaces and returns a pointer to a RequestParams struct
func DevideIDEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"device_id": params[1].(string),
		},
	}
}

// --- SEARCH (2 params - [0] = auth, [1] = query{q, type, limit})
// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func SearchItemsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"q":     params[1].(string),
			"type":  params[2].(string),
			"limit": "1",
		},
	}
}
