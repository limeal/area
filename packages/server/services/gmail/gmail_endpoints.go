package gmail

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
)

// It returns a map of endpoints for the Gmail API
func GmailEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Actions
		"GetMailInfoEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/messages/${mailId}",
			Params:         GetMailInfoEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetProfileEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/profile",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllMailEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/messages",
			Params:         GetAllMailEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllDraftMailEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/drafts",
			Params:         GetAllMailEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListLabelsEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/labels",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200, 204},
		},
		"ListFiltersEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/settings/filters",
			Params:         BasicEndpointParams,
			ExpectedStatus: []int{200, 204},
		},
		// Reactions
		"SendMailEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/messages/send",
			Params:         SendMailEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateDraftMailEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/drafts",
			Params:         CreateDraftMailEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateLabelEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/labels",
			Params:         CreateEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateFilterEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/settings/filters",
			Params:         CreateEndpointParams,
			ExpectedStatus: []int{200},
		},
		"UpdateVacationMessageEndpoint": {
			BaseURL:        "https://www.googleapis.com/gmail/v1/users/${id}/settings/vacation",
			Params:         UpdateVacationMessageEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It takes in an array of interfaces and returns a pointer to a RequestParams struct
func BasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
	}
}

// It returns a pointer to a `utils.RequestParams` struct with the `Method`, `Headers`, and `UrlParams`
// fields set
func GetMailInfoEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
		},
		UrlParams: map[string]string{
			"id":     params[1].(string),
			"mailId": params[2].(string),
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetAllMailEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
		QueryParams: params[2].(map[string]string),
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func SendMailEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
		Body: "{ \"raw\": \"" + params[2].(string) + "\"}",
	}
}

// It creates a request params object for the `CreateDraftMail` endpoint
func CreateDraftMailEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
		Body: "{ \"message\": { \"raw\": \"" + params[2].(string) + "\"}}",
	}
}

// It creates a request params object for the CreateEndpoint function
func CreateEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
		Body: params[2].(string),
	}
}

// It takes three parameters, the first one is a pointer to a `models.Authorization` struct, the second
// one is a string, and the third one is a string. It returns a pointer to a `utils.RequestParams`
// struct
func UpdateVacationMessageEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
		Body: params[2].(string),
	}
}
