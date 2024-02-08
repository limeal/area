package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// `PollChoice` is a struct with a single field, `Title`, which is a string.
// @property {string} Title - The title of the poll choice.
type PollChoice struct {
	Title string `json:"title"`
}

// `PollChoices` is a struct with a field `Choices` of type `[]PollChoice`.
//
// The `[]` indicates that `PollChoice` is an array.
//
// The `PollChoice` type is defined in the same file as `PollChoices`:
//
// // Go
// type PollChoice struct {
// 	Text string `json:"text"`
// 	Votes int `json:"votes"`
// }
// @property {[]PollChoice} Choices - An array of PollChoice objects.
type PollChoices struct {
	Choices []PollChoice `json:"choices"`
}

// `PollError` is a struct with two fields, `Message` and `Error`, both of which are strings.
// @property {string} Message - The message to display to the user.
// @property {string} Error - The error message
type PollError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// OK

// It creates a poll
func createPoll(req static.AreaRequest) shared.AreaResponse {
	userID := req.AuthStore["user_id"].(string)
	if (*req.Store)["req:streamer:login"] != nil {
		streamer, _, err := req.Service.Endpoints["GetUserByLoginEndpoint"].Call([]interface{}{req.Authorization, (*req.Store)["req:streamer:login"]})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		userID = streamer["id"].(string)
	}

	pollChoices := PollChoices{}
	tab := strings.Split((*req.Store)["req:poll:choices"].(string), ",")
	for _, choice := range tab {
		pollChoices.Choices = append(pollChoices.Choices, PollChoice{Title: choice})
	}

	duration, err := strconv.Atoi((*req.Store)["req:poll:duration"].(string))
	if err != nil || duration < 15 || duration > 1800 {
		return shared.AreaResponse{Error: err}
	}

	title := utils.GenerateFinalComponent((*req.Store)["req:poll:title"].(string), req.ExternalData, []string{})
	body := make(map[string]interface{})
	body["broadcaster_id"] = userID
	body["title"] = title
	body["choices"] = pollChoices.Choices
	body["duration"] = duration

	if (*req.Store)["req:poll:points:per:vote"] != nil {
		body["channel_points_voting_enabled"] = true
		var err error
		body["channel_points_per_vote"], err = strconv.Atoi((*req.Store)["req:poll:points:per:vote"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	ebody, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, httpResp, errr := req.Service.Endpoints["CreatePollEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		string(ebody),
	})

	if httpResp != nil && httpResp.StatusCode == 403 {
		return shared.AreaResponse{Error: errors.New("You must be a partner or affiliate to create a poll")}
	}

	return shared.AreaResponse{
		Error: errr,
	}
}

// It creates a poll
func DescriptorForTwitchReactionCreatePoll() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_poll",
		Description: "Create a poll (Must be a partner or affiliate)",
		RequestStore: map[string]static.StoreElement{
			"req:poll:title": {
				Priority:    1,
				Description: "The title of the poll",
				Required:    true,
			},
			"req:poll:choices": {
				Priority:    2,
				Type:        "long_string",
				Description: "The choices of the poll (format: choice1,choice2,choice3)",
				Required:    true,
			},
			"req:poll:duration": {
				Priority:    3,
				Description: "The duration of the poll in seconds (default: 60, min: 15, max: 1800)",
				Required:    true,
			},
			"req:poll:points:per:vote": {
				Priority:    4,
				Description: "The amount of points, a user spend per additional vote (default: disabled, min: 1, max: 1000000)",
				Required:    false,
			},
			"req:streamer:login": {
				Description: "The streamer login, you want to create a poll (default: you)",
				Required:    false,
			},
		},
		Method: createPoll,
	}
}
