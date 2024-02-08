package reddit

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	commonr "area-server/services/common"
)

// `RedditValidators` returns a `static.ServiceValidator` that contains two validators, one for the
// `subreddit` parameter and one for the `link` parameter
func RedditValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:subreddit:name": SubRedditNameValidator,
		"req:link:url":       commonr.URLValidator,
	}
}

// If the value is a string, and the string is a valid subreddit name, then return true
func SubRedditNameValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	_, _, err := service.Endpoints["GetSubredditEndpoint"].CallEncode([]interface{}{
		authorization,
		value,
	})

	if err != nil {
		return false
	}

	return true
}
