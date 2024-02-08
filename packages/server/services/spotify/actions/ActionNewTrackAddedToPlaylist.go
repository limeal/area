package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
)

// A Track is a struct with a string Id, a string Name, a string HRef, a string PreviewUrl, an int
// Duration, and an Album.
// @property {string} Id - A unique identifier for the track.
// @property {string} Name - The name of the track.
// @property {string} HRef - A link to the Web API endpoint providing full details of the track.
// @property {string} PreviewUrl - A link to a 30 second preview (MP3 format) of the track.
// @property {int} Duration - The duration of the track in milliseconds.
// @property {Album} Album - The album on which the track appears. The album object includes a link in
// href to full information about the album.
type Track struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	HRef       string `json:"href"`
	PreviewUrl string `json:"preview_url"`
	Duration   int    `json:"duration_ms"`
	Album      Album  `json:"album"`
}

// A TrackItem is a struct that contains a Track and a string.
// @property {Track} Track - A simplified track object (described above)
// @property {string} AddedAt - The date and time the track was saved. Note that some very old
// playlists may return null in this field.
type TrackItem struct {
	Track   Track  `json:"track"`
	AddedAt string `json:"added_at"`
}

// Tracks is a struct that contains a slice of TrackItem and an int.
// @property {[]TrackItem} Items - An array of TrackItem objects.
// @property {int} Total - The total number of tracks in the playlist.
type Tracks struct {
	Items []TrackItem `json:"items"`
	Total int         `json:"total"`
}

// SpotifyPlaylist is a struct that contains a string, a Tracks type, and a string.
// @property {string} Name - The name of the playlist
// @property {Tracks} Tracks - A list of tracks in the playlist.
// @property {string} SnapshotId - The version identifier for the current playlist. You can use this
// value to determine if a playlist has changed since the last time you retrieved it.
type SpotifyPlaylist struct {
	Name       string `json:"name"`
	Tracks     Tracks `json:"tracks"`
	SnapshotId string `json:"snapshot_id"`
}

// It checks if a new track has been added to a playlist
func hasNewTrackAddedToPlaylist(req static.AreaRequest) shared.AreaResponse {
	encode, _, err := req.Service.Endpoints["FindPlaylistByIDEndpoint"].CallEncode([]interface{}{req.Authorization, (*req.Store)["req:playlist:id"]})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	playlist := SpotifyPlaylist{}
	if err := json.Unmarshal(encode, &playlist); err != nil {
		return shared.AreaResponse{Error: err}
	}

	if (*req.Store)["ctx:snapshot:id"] == nil {
		(*req.Store)["ctx:snapshot:id"] = playlist.SnapshotId
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:snapshot:id"] == playlist.SnapshotId {
		return shared.AreaResponse{Success: false}
	}

	(*req.Store)["ctx:snapshot:id"] = playlist.SnapshotId
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			// Playlist
			"spotify:playlist:id":          (*req.Store)["req:playlist:id"], // 37i9dQZF1DXcBWIGoYBM5M
			"spotify:playlist:name":        playlist.Name,                   // Rock Classics
			"spotify:playlist:snapshot:id": playlist.SnapshotId,             // 0QJ4q7q7X3ZJZy9Y8X5Z0w
			// Track
			"spotify:track:uri":         "spotify:track:" + playlist.Tracks.Items[0].Track.Id, // spotify:track:6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:name":        playlist.Tracks.Items[0].Track.Name,                  // We Will Rock You
			"spotify:track:duration":    playlist.Tracks.Items[0].Track.Duration,              // 354000
			"spotify:track:href":        playlist.Tracks.Items[0].Track.HRef,                  // https://api.spotify.com/v1/tracks/6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:preview:url": playlist.Tracks.Items[0].Track.PreviewUrl,            // https://p.scdn.co/mp3-preview/08b7e...
			// Album
			"spotify:album:id":           "spotify:album:" + playlist.Tracks.Items[0].Track.Album.Id, // spotify:album:6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:name":         playlist.Tracks.Items[0].Track.Album.Name,                  // News Of The World
			"spotify:album:total_tracks": playlist.Tracks.Items[0].Track.Album.TotalTracks,           // 10
			"spotify:album:href":         playlist.Tracks.Items[0].Track.Album.HRef,                  // https://api.spotify.com/v1/albums/6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:release:date": playlist.Tracks.Items[0].Track.Album.ReleaseDate,           // 1977-11-10
			"spotify:album:type":         playlist.Tracks.Items[0].Track.Album.AlbumType,             // album
		},
	}
}

// It returns a `static.ServiceArea` that describes the Spotify action `new_track_added_to_playlist`
func DescriptorForSpotifyActionNewTrackAddedToPlaylist() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_track_added_to_playlist",
		Description: "Triggered when a playlist has a new track added to it",
		RequestStore: map[string]static.StoreElement{
			"req:playlist:id": {
				Description: "ID of the playlist to watch for new tracks",
				Required:    true,
				Type:        "select_uri",
				Values:      []string{"/playlists"},
			},
		},
		Method: hasNewTrackAddedToPlaylist,
		Components: []string{
			"spotify:playlist:id",          // 37i9dQZF1DXcBWIGoYBM5M
			"spotify:playlist:name",        // Rock Classics
			"spotify:playlist:snapshot:id", // 0QJ4q7q7X3ZJZy9Y8X5Z0w
			// Track
			"spotify:track:uri",         // spotify:track:6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:name",        // We Will Rock You
			"spotify:track:duration",    // 354000
			"spotify:track:href",        // https://api.spotify.com/v1/tracks/6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:preview:url", // https://p.scdn.co/mp3-preview/08b7e...
			// Album
			"spotify:album:id",           // spotify:album:6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:name",         // News Of The World
			"spotify:album:total_tracks", // 10
			"spotify:album:href",         // https://api.spotify.com/v1/albums/6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:release:date", // 1977-11-10
			"spotify:album:type",
		},
	}
}
