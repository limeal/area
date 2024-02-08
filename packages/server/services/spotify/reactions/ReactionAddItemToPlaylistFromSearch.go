package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"errors"
)

// `Item` is a struct with a single field, `URI`, which is a string.
// @property {string} URI - The URI of the item.
type Item struct {
	URI string `json:"uri"`
}

// `SearchItemsResponse` is a struct that contains a struct that contains a slice of `Item`s.
// @property Tracks - A list of tracks that match the search query.
// @property Albums - A list of albums that match the query.
// @property Artists - A list of artists that match the query.
// @property Playlists - A list of playlists that match the search query.
// @property Shows - A list of shows that match the search query.
// @property Episodes - A list of episodes that match the search query.
// @property AudioBooks - A list of audio books that match the search query.
type SearchItemsResponse struct {
	Tracks struct {
		Items []Item `json:"items"`
	} `json:"tracks"`
	Albums struct {
		Items []Item `json:"items"`
	} `json:"albums"`
	Artists struct {
		Items []Item `json:"items"`
	} `json:"artists"`
	Playlists struct {
		Items []Item `json:"items"`
	} `json:"playlists"`
	Shows struct {
		Items []Item `json:"items"`
	} `json:"shows"`
	Episodes struct {
		Items []Item `json:"items"`
	} `json:"episodes"`
	AudioBooks struct {
		Items []Item `json:"items"`
	} `json:"audiobooks"`
}

// It searches for an item, and adds the first result to a playlist
func addItemToPlaylistFromSearch(req static.AreaRequest) shared.AreaResponse {
	if req.ExternalData["spotify:playlist:id"] != nil && req.ExternalData["spotify:playlist:id"] == (*req.Store)["req:writable:playlist:id"] {
		return shared.AreaResponse{Error: errors.New("Cannot add to the same playlist")}
	}

	query := utils.GenerateFinalComponent((*req.Store)["req:search:query"].(string), req.ExternalData, []string{})
	playlistID := utils.GenerateFinalComponent((*req.Store)["req:writable:playlist:id"].(string), req.ExternalData, []string{
		"spotify:playlist:id",
	})

	// Search for the item
	typeItem := "track"
	if (*req.Store)["req:search:type"] != nil {
		typeItem = (*req.Store)["req:search:type"].(string)
	}

	encode, _, err := req.Service.Endpoints["SearchItemsEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		query,
		typeItem,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	searchItems := SearchItemsResponse{}
	if err := json.Unmarshal(encode, &searchItems); err != nil {
		return shared.AreaResponse{Error: err}
	}

	// Choose the first item of the correct type
	var items []Item
	switch typeItem {
	case "track":
		items = searchItems.Tracks.Items
	case "album":
		items = searchItems.Albums.Items
	case "artist":
		items = searchItems.Artists.Items
	case "playlist":
		items = searchItems.Playlists.Items
	case "show":
		items = searchItems.Shows.Items
	case "episode":
		items = searchItems.Episodes.Items
	case "audiobook":
		items = searchItems.AudioBooks.Items
	default:
		return shared.AreaResponse{Error: errors.New("Unknown type")}
	}

	if len(items) == 0 {
		return shared.AreaResponse{Error: errors.New("No item found")}
	}

	body := map[string]interface{}{
		"uris":     []interface{}{items[0].URI},
		"position": 0,
	}

	str, errr := json.Marshal(body)
	if errr != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr0 := req.Service.Endpoints["AddItemToPlaylistEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		playlistID,
		string(str),
	})

	if errr0 != nil {
		return shared.AreaResponse{Error: errr}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea that describes the service area
func DescriptorForSpotifyReactionAddItemToPlaylistFromSearch() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_item_to_playlist_from_search",
		Description: "Add an item to a playlist from a search",
		RequestStore: map[string]static.StoreElement{
			"req:writable:playlist:id": {
				Type:        "select_uri",
				Description: "The ID of the playlist",
				Required:    true,
				Values:      []string{"/playlists"},
			},
			"req:search:query": {
				Priority:    1,
				Type:        "string",
				Description: "The query to search",
				Required:    true,
			},
			"req:search:type": {
				Priority:    2,
				Type:        "select",
				Description: "The type of the search (default: track)",
				Required:    false,
				Values: []string{
					"track",
					"album",
					"artist",
					"playlist",
					"show",
					"episode",
					"audiobook",
				},
			},
		},
		Method: addItemToPlaylistFromSearch,
	}
}
