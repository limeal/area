package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// "Show is a struct with four fields, all of which are exported."
//
// The first line of the type definition is the type keyword, followed by the name of the type,
// followed by the keyword struct.
//
// The next four lines are the fields of the struct. Each field is a name, followed by a type, followed
// by a tag. The tag is a string literal that is used by the encoding/json package to encode and decode
// the struct to and from JSON.
//
// The last line of the type definition is a closing curly brace.
//
// The fields of the struct
// @property {string} Id - The unique identifier for the show.
// @property {string} Name - The name of the show
// @property {string} HRef - The URL to the show's page on the site.
// @property {int} TotalEpisodes - The total number of episodes in the show.
type Show struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	HRef          string `json:"href"`
	TotalEpisodes int    `json:"total_episodes"`
}

// `ShowItem` is a struct that contains a `Show` struct and a `string`
// @property {Show} Show - This is the show object that contains all the information about the show.
// @property {string} AddedAt - The date the show was added to the user's list.
type ShowItem struct {
	Show    Show   `json:"show"`
	AddedAt string `json:"added_at"`
}

// `Shows` is a struct that contains a slice of `ShowItem`s and an `int`
// @property {[]ShowItem} Items - An array of ShowItem objects.
// @property {int} Total - The total number of shows in the response.
type Shows struct {
	Items []ShowItem `json:"items"`
	Total int        `json:"total"`
}

// It checks if the user has a new saved show
func hasANewSavedShow(req static.AreaRequest) shared.AreaResponse {

	encode, _, err := req.Service.Endpoints["FindUserSavedShowsEndpoint"].CallEncode([]interface{}{req.Authorization})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	items := Shows{}
	if err := json.Unmarshal(encode, &items); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbItems := len(items.Items)
	ok, err := utils.IsLatestByDate(req.Store, nbItems, func() interface{} {
		return items.Items[0].Show.Id
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
			"spotify:show:uri":            "spotify:show:" + items.Items[0].Show.Id, // spotify:show:6rqhFgbbKwnb9MLmUQDhG6
			"spotify:show:name":           items.Items[0].Show.Name,                 // We Will Rock You
			"spotify:show:total_episodes": items.Items[0].Show.TotalEpisodes,
			"spotify:show:href":           items.Items[0].Show.HRef,
		},
	}
}

// `DescriptorForSpotifyActionNewSavedShow` returns a `static.ServiceArea` that describes the
// `new_saved_show` action
func DescriptorForSpotifyActionNewSavedShow() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_saved_show",
		Description: "When a new show is saved",
		Method:      hasANewSavedShow,
		Components: []string{
			"spotify:show:uri",
			"spotify:show:name",
			"spotify:show:total_episodes",
			"spotify:show:href",
		},
	}
}
