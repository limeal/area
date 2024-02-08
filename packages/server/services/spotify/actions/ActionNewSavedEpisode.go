package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// Episode is a struct that contains an id, name, and href.
// @property {string} Id - The id of the episode.
// @property {string} Name - The name of the episode
// @property {string} HRef - The URL to the episode's page on the TVDB website.
type Episode struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	HRef string `json:"href"`
}

// `EpisodeItem` is a struct that contains an `Episode` and a `string`
// @property {Episode} Episode - The episode object
// @property {string} AddedAt - The date and time the episode was added to the user's library.
type EpisodeItem struct {
	Episode Episode `json:"episode"`
	AddedAt string  `json:"added_at"`
}

// Episodes is a struct that contains a slice of EpisodeItem structs and an int.
// @property {[]EpisodeItem} Items - An array of EpisodeItem objects.
// @property {int} Total - The total number of episodes in the podcast.
type Episodes struct {
	Items []EpisodeItem `json:"items"`
	Total int           `json:"total"`
}

// It checks if the user has a new saved episode
func hasANewSavedEpisode(req static.AreaRequest) shared.AreaResponse {

	encode, _, err := req.Service.Endpoints["FindUserSavedEpisodesEndpoint"].CallEncode([]interface{}{req.Authorization})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	items := Episodes{}
	if err := json.Unmarshal(encode, &items); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbItems := len(items.Items)
	ok, err := utils.IsLatestByDate(req.Store, nbItems, func() interface{} {
		return items.Items[0].Episode.Id
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
			"spotify:episode:uri":  "spotify:episode:" + items.Items[0].Episode.Id,
			"spotify:episode:name": items.Items[0].Episode.Name,
		},
	}
}

// `DescriptorForSpotifyActionNewSavedEpisode` returns a `static.ServiceArea` that describes the
// `new_saved_episode` action for the `spotify` service
func DescriptorForSpotifyActionNewSavedEpisode() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_saved_episode",
		Description: "When a new episode is saved",
		Method:      hasANewSavedEpisode,
		Components: []string{
			"spotify:episode:uri",
			"spotify:episode:name",
		},
	}
}
