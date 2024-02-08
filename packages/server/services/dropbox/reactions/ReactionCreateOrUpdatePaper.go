package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"errors"
	"strings"
)

// It creates or updates a paper in Dropbox Paper
func createOrUpdatePaper(req static.AreaRequest) shared.AreaResponse {

	paperName := utils.GenerateFinalComponent((*req.Store)["req:paper:name"].(string), req.ExternalData, []string{})
	paperContent := utils.GenerateFinalComponent((*req.Store)["req:paper:content"].(string), req.ExternalData, []string{})

	if !strings.HasSuffix(paperName, ".paper") {
		return shared.AreaResponse{Error: errors.New("Paper name must end with .paper")}
	}

	if (*req.Store)["ctx:directory:path"] == nil {
		(*req.Store)["ctx:directory:path"] = ""
		if (*req.Store)["req:directory:path"] != nil {
			(*req.Store)["ctx:directory:path"] = (*req.Store)["req:directory:path"]
		}
	}

	if (*req.Store)["ctx:import:format"] == nil {
		(*req.Store)["ctx:import:format"] = "plain_text"
		if (*req.Store)["req:import:format"] != nil {
			(*req.Store)["ctx:import:format"] = (*req.Store)["req:import:format"]
		}
	}

	var rerr error
	var rbody map[string]interface{}
	path := (*req.Store)["ctx:directory:path"].(string) + "/" + paperName
	if (*req.Store)["ctx:paper:revision"] == nil {

		body, err := json.Marshal(map[string]interface{}{
			"path":          path,
			"import_format": (*req.Store)["ctx:import:format"],
		})

		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		rbody, _, rerr = req.Service.Endpoints["CreatePaperEndpoint"].Call([]interface{}{
			req.Authorization,
			string(body),
			paperContent,
		})
	} else {
		body, err := json.Marshal(map[string]interface{}{
			"path":              path,
			"import_format":     (*req.Store)["ctx:import:format"],
			"paper_revision":    (*req.Store)["ctx:paper:revision"],
			"doc_update_policy": "update",
		})

		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		rbody, _, rerr = req.Service.Endpoints["UpdatePaperEndpoint"].Call([]interface{}{
			req.Authorization,
			string(body),
			paperContent,
		})
	}

	if rerr != nil {
		return shared.AreaResponse{Error: rerr}
	}

	(*req.Store)["ctx:paper:revision"] = rbody["paper_revision"]
	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a `static.ServiceArea` object that describes the service area
func DescriptorForDropboxReactionCreateOrUpdatePaper() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_or_update_paper",
		Description: "Create or update a paper.",
		RequestStore: map[string]static.StoreElement{
			"req:paper:name": {
				Description: "The name of the paper to create or update",
				Type:        "string",
				Required:    true,
			},
			"req:paper:content": {
				Priority:    1,
				Description: "The content of the paper to create or update",
				Type:        "long_string",
				Required:    true,
			},
			"req:directory:path": {
				Priority:    2,
				Description: "The directory path where to create or update the paper",
				Type:        "select_uri",
				Required:    false,
				Values:      []string{"/directories"},
			},
			"req:import:format": {
				Priority:    3,
				Description: "The import format of the paper to create or update (default: plain_text)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"markdown",
					"html",
					"plain_text",
				},
			},
		},
		Method: createOrUpdatePaper,
	}
}
