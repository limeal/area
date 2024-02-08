package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It adds a file to a directory
func addFileToDirectory(req static.AreaRequest) shared.AreaResponse {

	filename := utils.GenerateFinalComponent((*req.Store)["req:file:name"].(string), req.ExternalData, []string{})
	fileContent := utils.GenerateFinalComponent((*req.Store)["req:file:content"].(string), req.ExternalData, []string{})

	if (*req.Store)["ctx:directory:path"] == nil {
		(*req.Store)["ctx:directory:path"] = ""
		if (*req.Store)["req:directory:path"] != nil {
			(*req.Store)["ctx:directory:path"] = (*req.Store)["req:directory:path"]
		}
	}

	if (*req.Store)["ctx:mode"] == nil {
		(*req.Store)["ctx:mode"] = "add"
		if (*req.Store)["req:allow:overwrite"] != nil && (*req.Store)["req:allow:overwrite"].(string) == "true" {
			(*req.Store)["ctx:mode"] = "overwrite"
		}
	}

	if (*req.Store)["ctx:mute"] == nil {
		(*req.Store)["ctx:mute"] = false
		if (*req.Store)["req:notify:users"] != nil && (*req.Store)["req:notify:users"].(string) == "false" {
			(*req.Store)["ctx:mute"] = true
		}
	}

	if (*req.Store)["ctx:allow:auto:rename"] == nil {
		(*req.Store)["ctx:allow:auto:rename"] = false
		if (*req.Store)["req:allow:auto:rename"] != nil && (*req.Store)["req:allow:auto:rename"].(string) == "true" {
			(*req.Store)["ctx:allow:auto:rename"] = true
		}
	}

	dropboxArg, err := json.Marshal(map[string]interface{}{
		"path":            (*req.Store)["ctx:directory:path"].(string) + "/" + filename,
		"mode":            (*req.Store)["ctx:mode"],
		"autorename":      (*req.Store)["ctx:allow:auto:rename"],
		"mute":            (*req.Store)["ctx:mute"],
		"strict_conflict": false,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, err = req.Service.Endpoints["UploadFileEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		string(dropboxArg),
		fileContent,
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForDropboxReactionAddFileToDirectory() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_file_to_directory",
		Description: "Add a file to a directory",
		RequestStore: map[string]static.StoreElement{
			"req:file:name": {
				Description: "The file name to add",
				Type:        "string",
				Required:    true,
			},
			"req:file:content": {
				Priority:    1,
				Description: "The file content to add",
				Type:        "long_string",
				Required:    true,
			},
			"req:allow:overwrite": {
				Priority:    3,
				Description: "Allow overwriting the file if it already exists",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
			"req:allow:auto:rename": {
				Priority:    3,
				Description: "Allow auto renaming the file if it already exists",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
			"req:notify:users": {
				Priority:    3,
				Description: "Notify users when the file is added (default: true)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
			"req:directory:path": {
				Priority:    2,
				Description: "The directory path where to add the file",
				Type:        "select_uri",
				Required:    false,
				Values:      []string{"/directories"},
			},
		},
		Method: addFileToDirectory,
	}
}
