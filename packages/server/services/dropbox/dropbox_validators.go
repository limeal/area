package dropbox

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/common"
	"encoding/json"
)

// It returns a map of validators for the Dropbox API
func DropboxValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:directory:path":           DirectoryPathValidator,
		"req:file:id":                  FileIdValidator,
		"req:allow:recursive":          common.BoolValidator,
		"req:allow:overwrite":          common.BoolValidator,
		"req:allow:auto:rename":        common.BoolValidator,
		"req:include:non:downloadable": common.BoolValidator,
		"req:entity:type":              EntityTypeValidator,
		"req:notify:users":             common.BoolValidator,
		"req:user:email":               common.EmailValidator,
		"req:access:level":             AccessLevelValidator,
	}
}

// ------------------ Validators ------------------

// It checks if the value is a string, and if it is, it checks if the value is a valid directory path
func DirectoryPathValidator(
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

	if value.(string) == "/" {
		return true
	}

	_, _, err := service.Endpoints["ListFoldersEndpoint"].CallEncode([]interface{}{authorization, "{\"path\": \"" + value.(string) + "\"}"})
	if err != nil {
		return false
	}

	return true
}

// It takes a file path, and returns true if the file exists
func FileIdValidator(
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

	body, errJ := json.Marshal(map[string]interface{}{
		"path":                                value.(string),
		"include_media_info":                  false,
		"include_deleted":                     false,
		"include_has_explicit_shared_members": false,
	})
	if errJ != nil {
		return false
	}

	_, _, err := service.Endpoints["GetFileMetadataEndpoint"].CallEncode([]interface{}{authorization, string(body)})
	if err != nil {
		return false
	}

	return true
}

// "If the value is a string, and it's either `group`, `user`, or `invitee`, then return true."
//
// The `authorization` parameter is the authorization object that was passed to the `Authorize`
// function. The `service` parameter is the service object that was passed to the `Authorize` function.
// The `value` parameter is the value of the field that is being validated. The `store` parameter is a
// map that can be used to store data that can be used by other validators
func EntityTypeValidator(
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

	switch value.(string) {
	case "group", "user", "invitee":
		return true
	default:
		return false
	}
}

// "If the value is a string, and it's either 'viewer', 'editor', or 'owner', then return true."
//
// The `value` parameter is the value of the field being validated. The `store` parameter is a map of
// all the values in the request
func AccessLevelValidator(
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

	switch value.(string) {
	case "viewer", "editor", "owner":
		return true
	default:
		return false
	}
}
