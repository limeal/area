package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"strconv"
)

// It creates a new gist
func createNewGist(req static.AreaRequest) shared.AreaResponse {
	// Find the user the first time

	name := utils.GenerateFinalComponent((*req.Store)["req:gist:name"].(string), req.ExternalData, []string{})
	description := utils.GenerateFinalComponent((*req.Store)["req:gist:description"].(string), req.ExternalData, []string{})
	content := utils.GenerateFinalComponent((*req.Store)["req:gist:content"].(string), req.ExternalData, []string{})

	public := false
	if (*req.Store)["req:gist:public"] != nil {
		var err error
		public, err = strconv.ParseBool((*req.Store)["req:gist:public"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	gist := &common.CreateGist{
		Description: description,
		Public:      public,
		Files: map[string]common.GistFile{
			name: {
				Content: content,
			},
		},
	}

	str, err := json.Marshal(gist)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["CreateNewGistEndpoint"].Call([]interface{}{req.Authorization, string(str)})
	if errr != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a `static.ServiceArea` object that describes the service area
func DescriptorForGithubReactionCreateNewGist() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_gist",
		Description: "Create a new gist",
		RequestStore: map[string]static.StoreElement{
			"req:gist:name": {
				Description: "Name of the gist",
				Required:    true,
			},
			"req:gist:description": {
				Priority:    1,
				Description: "Description of the gist",
				Required:    true,
			},
			"req:gist:content": {
				Priority:    2,
				Description: "Content of the gist",
				Required:    true,
			},
			"req:gist:public": {
				Priority: 3,
				Required: false,
				Type:     "select",
				Values:   []string{"true", "false"},
			},
		},
		Method: createNewGist,
	}
}
