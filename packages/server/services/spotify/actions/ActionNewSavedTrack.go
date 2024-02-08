package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It checks if the user has a new saved track, and if so, it returns the track's information
func hasANewSavedTrack(req static.AreaRequest) shared.AreaResponse {

	encode, _, err := req.Service.Endpoints["FindUserSavedTracksEndpoint"].CallEncode([]interface{}{req.Authorization})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	items := Tracks{}
	if err := json.Unmarshal(encode, &items); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbItems := len(items.Items)
	ok, err := utils.IsLatestByDate(req.Store, nbItems, func() interface{} {
		return items.Items[0].Track.Id
	}, func() string {
		return items.Items[0].AddedAt
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			// Track
			"spotify:track:uri":         "spotify:track:" + items.Items[0].Track.Id, // spotify:track:6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:name":        items.Items[0].Track.Name,                  // We Will Rock You
			"spotify:track:duration":    items.Items[0].Track.Duration,              // 354000
			"spotify:track:href":        items.Items[0].Track.HRef,                  // https://api.spotify.com/v1/tracks/6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:preview:url": items.Items[0].Track.PreviewUrl,            // https://p.scdn.co/mp3-preview/08b7e...
			// Album
			"spotify:album:uri":          "spotify:album:" + items.Items[0].Track.Album.Id, // spotify:album:6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:name":         items.Items[0].Track.Album.Name,                  // News of the World
			"spotify:album:total_tracks": items.Items[0].Track.Album.TotalTracks,           // 14
			"spotify:album:href":         items.Items[0].Track.Album.HRef,                  // https://api.spotify.com/v1/albums/6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:release:date": items.Items[0].Track.Album.ReleaseDate,           // 1977-11-10
			"spotify:album:type":         items.Items[0].Track.Album.AlbumType,             // album
		},
	}
}

// `DescriptorForSpotifyActionNewSavedTrack` returns a `static.ServiceArea` that describes the
// `new_saved_track` action for the `spotify` service
func DescriptorForSpotifyActionNewSavedTrack() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_saved_track",
		Description: "When a new track is saved",
		Method:      hasANewSavedTrack,
		Components: []string{
			"spotify:track:uri",
			"spotify:track:name",
			"spotify:track:duration",
			"spotify:track:href",
			"spotify:track:preview:url",

			"spotify:album:uri",
			"spotify:album:name",
			"spotify:album:total_tracks",
			"spotify:album:href",
			"spotify:album:release:date",
			"spotify:album:type",
		},
	}
}
