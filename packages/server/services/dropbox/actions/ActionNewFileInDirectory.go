package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/dropbox/common"
	"area-server/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

// `ListFoldersContinueResponse` is a struct with three fields: `Entries` of type `[]common.Entry`,
// `Cursor` of type `string`, and `HasMore` of type `bool`.
//
// The `Entries` field is an array of `common.Entry`s. The `common.Entry` type is defined in the
// `common` package.
//
// The `Cursor` field is a string.
//
// The `HasMore` field is a boolean.
// @property {[]common.Entry} Entries - An array of entries that are contained in the folder.
// @property {string} Cursor - A cursor that can be used to retrieve the next batch of results.
// @property {bool} HasMore - If true, then there are more entries available. Pass the cursor to
// ListFolderContinue() to retrieve the rest.
type ListFoldersContinueResponse struct {
	Entries []common.Entry `json:"entries"`
	Cursor  string         `json:"cursor"`
	HasMore bool           `json:"has_more"`
}

// It gets the latest file in a directory
func onNewFileInDirectory(req static.AreaRequest) shared.AreaResponse {
	if (*req.Store)["ctx:directory:path"] == nil {
		(*req.Store)["ctx:directory:path"] = (*req.Store)["req:directory:path"]
		if (*req.Store)["ctx:directory:path"] == "/" {
			(*req.Store)["ctx:directory:path"] = ""
		}
	}

	(*req.Store)["ctx:allow:recursive"] = false
	if (*req.Store)["req:allow:recursive"] != nil {
		var err error
		(*req.Store)["ctx:allow:recursive"], err = strconv.ParseBool((*req.Store)["req:allow:recursive"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	(*req.Store)["ctx:include:non:downloadable"] = true
	if (*req.Store)["req:include:non:downloadable"] != nil {
		var err error
		(*req.Store)["ctx:include:non:downloadable"], err = strconv.ParseBool((*req.Store)["req:include:non:downloadable"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	getLastCursor := func() (string, error) {
		body, err := json.Marshal(map[string]interface{}{
			"path":                                (*req.Store)["ctx:directory:path"],
			"recursive":                           (*req.Store)["ctx:allow:recursive"],
			"include_media_info":                  false,
			"include_deleted":                     false,
			"include_has_explicit_shared_members": false,
			"include_mounted_folders":             true,
			"include_non_downloadable_files":      (*req.Store)["ctx:include:non:downloadable"],
		})

		if err != nil {
			return "", err
		}

		lcursor, _, err0 := req.Service.Endpoints["ListFoldersGetLatestCursorEndpoint"].Call([]interface{}{
			req.Authorization,
			string(body),
		})
		if err0 != nil {
			return "", err0
		}

		return lcursor["cursor"].(string), nil
	}

	if (*req.Store)["ctx:cursor"] == nil {
		var err error
		(*req.Store)["ctx:cursor"], err = getLastCursor()
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	elements, _, errr := req.Service.Endpoints["ListFoldersContinueEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:cursor"],
	})
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	var response ListFoldersContinueResponse
	err := json.Unmarshal(elements, &response)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	fmt.Println("Elements: ", response.Entries)
	nbElements := len(response.Entries)
	ok, errL := utils.IsLatestBasic(req.Store, nbElements)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	var latestElement *common.Entry
	// Latest element correspond to the first file in the list at the end
	for i := nbElements - 1; i >= 0; i-- {
		if response.Entries[i].Tag == "file" {
			latestElement = &response.Entries[i]
		}
	}

	(*req.Store)["ctx:cursor"], err = getLastCursor()
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if latestElement == nil {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"dropbox:file:name":        latestElement.Name,
			"dropbox:file:path":        latestElement.PathLower,
			"dropbox:file:id":          latestElement.ID,
			"dropbox:file:size":        latestElement.Size,
			"dropbox:file:rev":         latestElement.Rev,
			"dropbox:file:client:time": latestElement.ClientModified,
			"dropbox:file:server:time": latestElement.ServerModified,
		},
	}
}

// It returns a static.ServiceArea object that describes the action
func DescriptorForDropboxActionNewFileInDirectory() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_file_in_directory",
		Description: "Triggered when a new file is created in a directory",
		RequestStore: map[string]static.StoreElement{
			"req:directory:path": {
				Description: "The path to the directory",
				Type:        "select_uri",
				Required:    true,
				Values:      []string{"/directories"},
			},
			"req:allow:recursive": {
				Priority:    1,
				Description: "Check for new files recursively",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
			"req:include:non:downloadable": {
				Priority:    1,
				Description: "Include non downloadable files or not",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
		},
		Method: onNewFileInDirectory,
		Components: []string{
			"dropbox:file:name",
			"dropbox:file:path",
			"dropbox:file:id",
			"dropbox:file:size",
			"dropbox:file:rev",
			"dropbox:file:client:time",
			"dropbox:file:server:time",
		},
	}
}
