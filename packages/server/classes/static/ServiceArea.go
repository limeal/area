package static

import (
	"area-server/classes/shared"
	"area-server/db/postgres/models"
	"fmt"
)

// `ServiceArea` is a struct that contains a name, description, a boolean value, a slice of strings, a
// map of `StoreElement`s, another boolean value, and a function that takes an `AreaRequest` and
// returns a `shared.AreaResponse`.
// @property {string} Name - The name of the service area. This is used to identify the area in the
// API.
// @property {string} Description - A short description of the service area
// @property {bool} UseGateway - If the service area is using the gateway, it will be able to use the
// gateway's components.
// @property {[]string} Components - A list of components that can be used to fill fields in the
// request.
// @property RequestStore - This is a map of the elements that are required to make the request.
// @property {bool} WIP - Is the area still in development
// @property Method - This is the function that will be called when the service area is requested.
type ServiceArea struct {
	Name         string                                  `json:"name"`
	Description  string                                  `json:"description"`
	UseGateway   bool                                    `json:"use_gateway"`
	Components   []string                                `json:"components"` // Components that can be used to fill fields
	RequestStore map[string]StoreElement                 `json:"store"`      // Elements that are required to make the request
	WIP          bool                                    `json:"wip"`        // Is the area still in development
	Method       (func(AreaRequest) shared.AreaResponse) `json:"-"`
}

type StoreElement struct {
	Priority          int      `json:"priority"`           // Priority of the field (0 = highest)
	Type              string   `json:"type"`               // (select_uri, bool, string, number, ...)
	Description       string   `json:"description"`        // Description of the field
	Required          bool     `json:"required"`           // Is the field required
	NeedFields        []string `json:"need_fields"`        // Fields that are required to show this field
	AllowedComponents []string `json:"allowed_components"` // Components that are allowed to be used for this field
	Values            []string `json:"values"`             // For select_uri, values[0] is the endpoint. For select, values are the values that can be selected
}

// MergeStore takes two maps and returns a new map that contains all the keys and values from both
// maps.
func MergeStore(bs map[string]StoreElement, rs map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	// Merge the two maps

	for key, value := range rs {
		result[key] = value
	}
	return result
}

// Validating the request.
func (a *ServiceArea) Validate(
	authorization *models.Authorization,
	service *Service,
	store map[string]interface{},
) map[string]bool {
	result := make(map[string]bool)
	ok := true

	fmt.Println("Store: ", store)
	for key, value := range a.RequestStore {
		if value.Required && store[key] == nil {
			result[key] = false
			ok = false
		}
	}

	for key, value := range store {
		if service.Validators[key] != nil {
			result[key] = service.Validators[key](authorization, service, value, store)
			if !result[key] {
				ok = false
			}
		}
	}

	if !ok {
		return result
	}

	return nil
}
