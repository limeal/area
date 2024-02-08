package spotify

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"strings"
)

// It returns a map of validators that can be used to validate the request parameters of the Spotify
// API
func SpotifyValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:device:id":            DeviceIDValidator,
		"req:uri":                  UriValidator,
		"req:uris":                 UrisValidator,
		"req:playlist:id":          PlaylistIDValidator,
		"req:writable:playlist:id": WriteablePlaylistIDValidator,
		"req:search:type":          SearchTypeValidator,
	}
}

// ------------------------- Validators ------------------------------

// It checks if the value is a valid device ID for the user
func DeviceIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	body, _, err := service.Endpoints["GetUserAvailableDevicesEndpoint"].Call([]interface{}{authorization})
	if err != nil || body["devices"] == nil {
		return false
	}

	for _, device := range body["devices"].([]interface{}) {
		if device.(map[string]interface{})["id"] == value {
			return true
		}
	}
	return false
}

// It checks if the value is a string, and if it is, it checks if it's a valid Spotify URI
func UriValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	uri := strings.Split(value.(string), ":")

	if len(uri) != 3 {
		return false
	}

	if uri[0] == "{{spotify" && uri[2] == "uri}}" {
		return true
	}

	if uri[0] != "spotify" {
		return false
	}

	endpoint := "Find" + strings.Title(uri[1]) + "ByIDEndpoint"
	if service.Endpoints[endpoint] == nil {
		return false
	}

	_, _, err := service.Endpoints[endpoint].CallEncode([]interface{}{authorization, uri[2]})
	if err != nil {
		return false
	}
	return true
}

// `PlaylistIDValidator` takes an authorization, a service, a value, and a store, and returns a boolean
func PlaylistIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{spotify:playlist:id}}" {
		return true
	}

	_, _, err := service.Endpoints["FindPlaylistByIDEndpoint"].CallEncode([]interface{}{authorization, value})
	if err != nil {
		return false
	}
	return true
}

// It checks that the value is a string and that it is a valid URI
func UrisValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	for _, uri := range value.([]interface{}) {
		if !UriValidator(authorization, service, uri, store) {
			return false
		}
	}
	return true
}

// If the value is nil, or not a string, or the string is not the special value, then we call the
// FindPlaylistByIDEndpoint to get the playlist. If the playlist is collaborative, or the owner of the
// playlist is the user, then the value is valid
func WriteablePlaylistIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{spotify:playlist:id}}" {
		return true
	}

	body, _, err := service.Endpoints["FindPlaylistByIDEndpoint"].Call([]interface{}{authorization, value})
	if err != nil {
		return false
	}

	if body["collaborative"].(bool) {
		return true
	}

	userBody, _, errP := service.Endpoints["GetUserProfileEndpoint"].Call([]interface{}{authorization})
	if errP != nil {
		return false
	}

	if body["owner"].(map[string]interface{})["id"] != userBody["id"] {
		return false
	}
	return true
}

// If the value is a string, and the string is one of the following: album, artist, playlist, track,
// show, episode, audiobook, then return true
func SearchTypeValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	switch value.(string) {
	case "album", "artist", "playlist", "track", "show", "episode", "audiobook":
		return true
	}
	return false
}
