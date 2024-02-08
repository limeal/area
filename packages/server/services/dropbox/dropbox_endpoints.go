package dropbox

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
)

// It returns a map of all the endpoints that Dropbox has, and the parameters that each endpoint
// expects
func DropboxEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Validators
		"GetFileMetadataEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/get_metadata",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListFoldersEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/list_folder",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"ListFoldersContinueEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/list_folder/continue",
			Params:         ListFoldersContinueEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListFoldersGetLatestCursorEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/list_folder/get_latest_cursor",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListFileMembersEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/sharing/list_file_members",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetSharedLinkFileEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/sharing/get_shared_link_file",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"ShareFileToUserEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/sharing/add_file_member",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateFolderEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/create_folder_v2",
			Params:         BasicPostEndpointParams,
			ExpectedStatus: []int{200},
		},
		"UploadFileEndpoint": {
			BaseURL:        "https://content.dropboxapi.com/2/files/upload",
			Params:         UploadFileEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreatePaperEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/paper/create",
			Params:         UploadFileEndpointParams,
			ExpectedStatus: []int{200},
		},
		"UpdatePaperEndpoint": {
			BaseURL:        "https://api.dropboxapi.com/2/files/paper/update",
			Params:         UploadFileEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

// --------------------------- Parameters ---------------------------

// It takes a slice of interfaces and returns a pointer to a RequestParams struct
func BasicPostEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		Body: params[1].(string),
	}
}

// It takes a list of parameters, and returns a pointer to a `RequestParams` struct
func ListFoldersContinueEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":  "application/json",
		},
		Body: "{\"cursor\": \"" + params[1].(string) + "\"}",
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func UploadFileEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization":   "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Content-Type":    "application/octet-stream",
			"Dropbox-API-Arg": params[1].(string),
			"Accept":          "application/json",
		},
		Body: params[2].(string),
	}
}
