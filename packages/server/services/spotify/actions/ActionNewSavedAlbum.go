package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// `Album` is a type that represents a Spotify album.
// @property {string} Id - The Spotify ID for the album.
// @property {string} Name - The name of the album.
// @property {int} TotalTracks - The total number of tracks in the album.
// @property {string} HRef - The Spotify Web API endpoint providing full details of the album.
// @property {string} ReleaseDate - The date the album was first released, for example 1981. Depending
// on the precision, it might be shown as 1981-12 or 1981-12-15.
// @property {string} AlbumType - The type of the album: one of "album", "single", or "compilation".
type Album struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	TotalTracks int    `json:"total_tracks"`
	HRef        string `json:"href"`
	ReleaseDate string `json:"release_date"`
	AlbumType   string `json:"album_type"`
}

// `AlbumItem` is a struct that contains an `Album` and a `string`.
// @property {Album} Album - The album object.
// @property {string} AddedAt - The date and time the album was saved. Note that some very old albums
// may not have this value.
type AlbumItem struct {
	Album   Album  `json:"album"`
	AddedAt string `json:"added_at"`
}

// `Albums` is a struct that contains a slice of `AlbumItem`s and an integer.
// @property {[]AlbumItem} Items - An array of AlbumItem objects.
// @property {int} Total - The total number of albums in the user's library.
type Albums struct {
	Items []AlbumItem `json:"items"`
	Total int         `json:"total"`
}

// It checks if the user has a new saved album
func hasANewSavedAlbum(req static.AreaRequest) shared.AreaResponse {
	encode, _, err := req.Service.Endpoints["FindUserSavedAlbumsEndpoint"].CallEncode([]interface{}{req.Authorization})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	items := Albums{}
	if err := json.Unmarshal(encode, &items); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbItems := len(items.Items)
	ok, err := utils.IsLatestByDate(req.Store, nbItems, func() interface{} {
		return items.Items[0].Album.Id
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
			"spotify:album:uri":          "spotify:album:" + items.Items[0].Album.Id,
			"spotify:album:name":         items.Items[0].Album.Name,
			"spotify:album:total_tracks": items.Items[0].Album.TotalTracks,
			"spotify:album:release:date": items.Items[0].Album.ReleaseDate,
			"spotify:album:type":         items.Items[0].Album.AlbumType,
		},
	}
}

// `DescriptorForSpotifyActionNewSavedAlbum` returns a `static.ServiceArea` that describes the
// `new_saved_album` action
func DescriptorForSpotifyActionNewSavedAlbum() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_saved_album",
		Description: "When a new album is saved",
		Method:      hasANewSavedAlbum,
		Components: []string{
			"spotify:album:uri",
			"spotify:album:name",
			"spotify:album:total_tracks",
			"spotify:album:release:date",
			"spotify:album:type",
		},
	}
}
