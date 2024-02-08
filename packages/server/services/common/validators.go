package common

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"regexp"
	"strconv"
)

// If the value is nil, return false. If the value is "true" or "false", return true. Otherwise, return
// false
func BoolValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if value == "true" || value == "false" {
		return true
	}

	return false
}

// If the value is not nil, not an empty string, and can be converted to an integer, then it's valid
func IntValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if value.(string) == "" {
		return false
	}

	if v, err := strconv.Atoi(value.(string)); err != nil || v < 0 {
		return false
	}

	return true
}

// If the value is not nil, not an empty string, and can be parsed as a float, then return true
func FloatValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if value.(string) == "" {
		return false
	}

	if v, err := strconv.ParseFloat(value.(string), 64); err != nil || v < 0 {
		return false
	}

	return true
}

// It takes the value of the field, and checks if it's a valid email address
func EmailValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	body, _, _, err := utils.DoRequest("https://isitarealemail.com/api/email/validate", &utils.RequestParams{
		Method: "GET",
		QueryParams: map[string]string{
			"email": value.(string),
		},
	}, []int{200}, true)

	if err != nil || body.(map[string]interface{})["status"] == nil || body.(map[string]interface{})["status"].(string) != "valid" {
		return false
	}

	return true
}

// It checks if the value is a string and if it's a valid URL
func URLValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	// Regex to check if the url is valid (http or https allowed)
	regex := regexp.MustCompile(`^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	if !regex.MatchString(value.(string)) {
		return false
	}

	return true
}
