package twitch

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/twitch/common"
	"encoding/json"
	"strconv"
	"strings"
)

// It returns a map of validators that can be used to validate the request parameters of the Twitch API
func TwitchValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:streamer:login":        UserLoginValidator,
		"req:user:login":            UserLoginValidator,
		"req:clip:entity:type":      ClipEntityTypeValidator,
		"req:clip:entity:value":     ClipEntityValidator,
		"req:game:name":             GameNameValidator,
		"req:poll:title":            PollTitleValidator,
		"req:poll:choices":          PollChoicesValidator,
		"req:poll:duration":         PollDurationValidator,
		"req:poll:points:per:vote":  PollPointsPerVoteValidator,
		"req:streamer:entity:type":  StreamerEntityTypeValidator,
		"req:streamer:entity:value": StreamerEntityValueValidator,
	}
}

// ----------------------- Validators -----------------------

// `UserLoginValidator` is a function that takes in an `Authorization` object, a `Service` object, a
// value, and a map of values. It returns a boolean
func UserLoginValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{twitch:user:login}}" {
		return true
	}

	encode, _, err := service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{auth, value.(string)})
	if err != nil {
		return false
	}

	streamer := common.TwitchUsers{}
	if err := json.Unmarshal(encode, &streamer); err != nil {
		return false
	}

	if len(streamer.Data) == 0 {
		return false
	}

	return true
}

// "If the value is a string, and it's either 'game' or 'user', then it's valid."
//
// The first thing we do is check if the value is nil. If it is, then we return false. This is because
// we don't want to allow nil values
func ClipEntityTypeValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	switch value.(string) {
	case "game", "user":
		return true
	}

	return false
}

// "If the value is a string, and the entity type is either 'game' or 'user', then validate the value
// as a game name or user login, respectively."
//
// The first thing we do is check if the value is nil. If it is, we return false. This is because we
// don't want to validate a nil value
func ClipEntityValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:entity:type"] == nil {
		return false
	}

	switch store["req:entity:type"].(string) {
	case "game":
		if GameNameValidator(auth, service, value, store) {
			return true
		}
	case "user":
		if UserLoginValidator(auth, service, value, store) {
			return true
		}
	}

	return false
}

// It checks if the value is a string, and if it is, it checks if the value is a valid game name
func GameNameValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value.(string) == "{{twitch:game:name}}" {
		return true
	}

	encode, _, err := service.Endpoints["GetGameEndpoint"].CallEncode([]interface{}{auth, value.(string)})
	if err != nil {
		return false
	}

	game := common.TwitchGames{}
	if err := json.Unmarshal(encode, &game); err != nil {
		return false
	}

	if len(game.Data) == 0 {
		return false
	}

	return true
}

// If the value is a string and it's less than 60 characters, then it's valid
func PollTitleValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if len(value.(string)) > 60 {
		return false
	}

	return true
}

// It checks that the value is a string, that it contains at least two and at most ten comma-separated
// values, and that each value is at most 25 characters long
func PollChoicesValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	tab := strings.Split(value.(string), ",")

	if len(tab) < 2 || len(tab) > 10 {
		return false
	}

	for _, choice := range tab {
		if len(choice) > 25 {
			return false
		}
	}

	return true
}

// If the value is a string, convert it to an integer and make sure it's between 15 and 1800
func PollDurationValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	duration, err := strconv.Atoi(value.(string))
	if err != nil {
		return false
	}

	if duration < 15 || duration > 1800 {
		return false
	}

	return true
}

// "If the value is a string, convert it to an integer and make sure it's between 1 and 1,000,000."
//
// The first thing we do is check if the value is nil. If it is, we return false
func PollPointsPerVoteValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	points, err := strconv.Atoi(value.(string))
	if err != nil {
		return false
	}

	if points < 1 || points > 1000000 {
		return false
	}

	return true
}

// `StreamerEntityTypeValidator` is a function that takes an authorization, a service, a value, and a
// store, and returns a boolean
func StreamerEntityTypeValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	switch value.(string) {
	case "id", "login":
		return true
	}

	return false
}

// It checks if the value is a string, and if it is, it checks if the value is a valid streamer entity
func StreamerEntityValueValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:streamer:entity:type"] == nil {
		return false
	}

	var encode []byte
	var err error
	switch store["req:streamer:entity:type"].(string) {
	case "id":

		if value.(string) == "{{twitch:user:id}}" || value.(string) == "{{twitch:creator:id}}" || value.(string) == "{{twitch:broadcaster:id}}" {
			return true
		}

		encode, _, err = service.Endpoints["GetUserByIdEndpoint"].CallEncode([]interface{}{auth, value.(string)})
	case "login":

		if value.(string) == "{{twitch:user:login}}" {
			return true
		}

		encode, _, err = service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{auth, value.(string)})
	}

	if err != nil {
		return false
	}

	streamer := common.TwitchUsers{}
	if err := json.Unmarshal(encode, &streamer); err != nil {
		return false
	}

	if len(streamer.Data) == 0 {
		return false
	}

	return true
}
