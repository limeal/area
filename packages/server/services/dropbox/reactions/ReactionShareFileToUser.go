package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It shares a file to a user
func shareFileToUser(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["ctx:custom:message"] == nil {
		(*req.Store)["ctx:custom:message"] = "Area - You have been shared a file."
		if (*req.Store)["req:custom:message"] != nil {
			(*req.Store)["ctx:custom:message"] = utils.GenerateFinalComponent((*req.Store)["req:custom:message"].(string), req.ExternalData, []string{})
		}
	}

	if (*req.Store)["ctx:access:level"] == nil {
		(*req.Store)["ctx:access:level"] = "viewer"
		if (*req.Store)["req:access:level"] != nil {
			(*req.Store)["ctx:access:level"] = (*req.Store)["req:access:level"]
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"file": (*req.Store)["req:file:id"],
		"members": []map[string]interface{}{
			{
				".tag":  "email",
				"email": (*req.Store)["req:user:email"],
			},
		},
		"custom_message": (*req.Store)["ctx:custom:message"],
		"access_level":   (*req.Store)["ctx:access:level"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, err = req.Service.Endpoints["ShareFileToUserEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		string(body),
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForDropboxReactionShareFileToUser() static.ServiceArea {
	return static.ServiceArea{
		Name:        "share_file_to_user",
		Description: "Share a file to a user.",
		RequestStore: map[string]static.StoreElement{
			"req:file:id": {
				Description: "The file to share",
				Type:        "select_uri",
				Required:    true,
				Values:      []string{"/files?path=${req:directory:path}"},
			},
			"req:user:email": {
				Priority:    1,
				Description: "The email of the user to share the file to",
				Type:        "string",
				Required:    true,
			},
			"req:custom:message": {
				Priority:    3,
				Description: "A custom message to include in the notification email (default: Area - You have been shared a file.)",
				Type:        "long_string",
				Required:    false,
			},
			"req:access:level": {
				Priority:    3,
				Description: "The access level of the user (default: viewer)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"viewer",
					"editor",
					"owner",
				},
			},
			"req:directory:path": {
				Priority:    2,
				Description: "The path to the directory",
				Type:        "select_uri",
				Required:    false,
				Values:      []string{"/directories"},
			},
		},
		Method: shareFileToUser,
	}
}
