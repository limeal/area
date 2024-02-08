package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/twitch/common"
	"encoding/json"
	"errors"
)

// It checks if the current track has changed
func hasCurrentTrackChange(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["ctx:streamer:id"] == nil && ((*req.Store)["req:streamer:login"] != nil && (*req.Store)["req:streamer:login"] != "") {
		encode, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:streamer:login"]})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		user := common.TwitchUsers{}
		if err := json.Unmarshal(encode, &user); err != nil {
			return shared.AreaResponse{Error: err}
		}

		if len(user.Data) == 0 {
			return shared.AreaResponse{Error: errors.New("Streamer not found")}
		}

		(*req.Store)["ctx:streamer:id"] = user.Data[0].ID
	}

	encode, httpResp, err := req.Service.Endpoints["GetSoundTrackCurrentTrackEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:streamer:id"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if httpResp.StatusCode == 404 { // Not playing a track
		return shared.AreaResponse{Success: false}
	}

	response := common.TwitchSoundTracks{}
	if err := json.Unmarshal(encode, &response.Data); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbTracks := len(response.Data)
	if (*req.Store)["ctx:tracks:total"] == nil {
		(*req.Store)["ctx:tracks:total"] = nbTracks
	}

	if (*req.Store)["ctx:tracks:total"] == 0 && nbTracks == 0 {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:last_track"] == nil {
		(*req.Store)["ctx:last_track"] = response.Data[0].Track.ID
	}

	if (*req.Store)["ctx:last_track"] == response.Data[0].Track.ID {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"twitch:track:title":               response.Data[0].Track.Title,
			"twitch:track:artist:id":           response.Data[0].Track.Artists[0].ID,
			"twitch:track:artist:name":         response.Data[0].Track.Artists[0].Name,
			"twitch:track:album:id":            response.Data[0].Track.Album.ID,
			"twitch:track:album:name":          response.Data[0].Track.Album.Name,
			"twitch:track:album:image":         response.Data[0].Track.Album.ImageURL,
			"twitch:track:duration":            response.Data[0].Track.Duration,
			"twitch:track:isrc":                response.Data[0].Track.ISRC,
			"twitch:track:id":                  response.Data[0].Track.ID,
			"twitch:track:source:id":           response.Data[0].Source.ID,
			"twitch:track:source:content_type": response.Data[0].Source.ContentType,
			"twitch:track:source:title":        response.Data[0].Source.Title,
			"twitch:track:source:image":        response.Data[0].Source.ImageURL,
			"twitch:track:source:soundtrack":   response.Data[0].Source.SoundtrackURL,
			"twitch:track:source:spotify":      response.Data[0].Source.SpotifyURL,
		},
	}
}

// It returns a `static.ServiceArea` object that describes the action
func DescriptorForTwitchActionCurrentTrackChange() static.ServiceArea {
	return static.ServiceArea{
		Name:        "current_track_change",
		Description: "When the current track change on your stream",
		WIP:         true,
		RequestStore: map[string]static.StoreElement{
			"req:streamer:login": {
				Description: "The streamer login, you want to check the current track (default: you)",
				Required:    false,
			},
		},
		Method: hasCurrentTrackChange,
		Components: []string{
			"twitch:track:title",
			"twitch:track:artist:id",
			"twitch:track:artist:name",
			"twitch:track:album:id",
			"twitch:track:album:name",
			"twitch:track:album:image",
			"twitch:track:duration",
			"twitch:track:isrc",
			"twitch:track:id",
			"twitch:track:source:id",
			"twitch:track:source:content_type",
			"twitch:track:source:title",
			"twitch:track:source:image",
			"twitch:track:source:soundtrack",
			"twitch:track:source:spotify",
		},
	}
}
