package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"regexp"
)

// `ListFiltersResponse` is a struct that contains a slice of structs that contain a struct that
// contains a struct that contains a slice of strings and a string and an int.
// @property {[]struct {
// 		ID       string `json:"id"`
// 		Criteria struct {
// 			From         string `json:"from"`
// 			To           string `json:"to"`
// 			Subject      string `json:"subject"`
// 			Query        string `json:"query"`
// 			NegatedQuery string `json:"negatedQuery"`
// 			Size         int    `json:"size"`
// 		} `json:"criteria"`
// 		Action struct {
// 			AddLabelIds    []string `json:"addLabelIds"`
// 			RemoveLabelIds []string `json:"removeLabelIds"`
// 			Forward        string   `json:"forward"`
// 		} `json:"action"`
// 	}} Filter - This is the array of filters that are returned.
type ListFiltersResponse struct {
	Filter []struct {
		ID       string `json:"id"`
		Criteria struct {
			From         string `json:"from"`
			To           string `json:"to"`
			Subject      string `json:"subject"`
			Query        string `json:"query"`
			NegatedQuery string `json:"negatedQuery"`
			Size         int    `json:"size"`
		} `json:"criteria"`
		Action struct {
			AddLabelIds    []string `json:"addLabelIds"`
			RemoveLabelIds []string `json:"removeLabelIds"`
			Forward        string   `json:"forward"`
		} `json:"action"`
	} `json:"messages"`
}

// It checks if the latest filter created matches the regex provided in the request, and if so, it
// returns the filter's data
func onNewFilterCreated(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["req:user:id"] != nil {
		(*req.Store)["ctx:user:id"] = (*req.Store)["req:user:id"]
	} else {
		(*req.Store)["ctx:user:id"] = "me"
	}

	encode, _, err := req.Service.Endpoints["ListFiltersEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var response ListFiltersResponse
	if err := json.Unmarshal(encode, &response); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbFilters := len(response.Filter)
	if ok, errr := utils.IsLatestBasic(req.Store, nbFilters); errr != nil || !ok {
		return shared.AreaResponse{Error: errr, Success: false}
	}

	newFilter := response.Filter[nbFilters-1]

	if (*req.Store)["req:filter:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:filter:regex"].(string), newFilter.ID)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		if !match {
			return shared.AreaResponse{Success: false}
		}

	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"gmail:user:id":             (*req.Store)["ctx:user:id"],
			"gmail:filter:id":           newFilter.ID,
			"gmail:filter:from":         newFilter.Criteria.From,
			"gmail:filter:to":           newFilter.Criteria.To,
			"gmail:filter:subject":      newFilter.Criteria.Subject,
			"gmail:filter:query":        newFilter.Criteria.Query,
			"gmail:filter:negatedQuery": newFilter.Criteria.NegatedQuery,
			"gmail:filter:size":         newFilter.Criteria.Size,
			"gmail:filter:forward":      newFilter.Action.Forward,
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForGmailActionNewFilterCreated() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_filter_created",
		Description: "When a new filter is created",
		WIP:         true,
		RequestStore: map[string]static.StoreElement{
			"req:filter:regex": {
				Priority:    1,
				Description: "The regex to check the filter for",
				Type:        "string",
				Required:    false,
			},
			"req:user:id": {
				Description: "The user id to check the mail for",
				Type:        "string",
				Required:    false,
			},
		},
		Method: onNewFilterCreated,
		Components: []string{
			"gmail:filter:id",
			"gmail:filter:from",
			"gmail:filter:to",
			"gmail:filter:subject",
			"gmail:filter:query",
			"gmail:filter:negatedQuery",
			"gmail:filter:size",
			"gmail:filter:forward",
		},
	}
}
