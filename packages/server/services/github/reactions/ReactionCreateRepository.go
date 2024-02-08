package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// It creates a new repository on GitHub
func createNewRepository(req static.AreaRequest) shared.AreaResponse {

	repoName := utils.GenerateFinalComponent((*req.Store)["req:new:repository:name"].(string), req.ExternalData, []string{})
	repoDescription := ""
	if (*req.Store)["req:new:repository:description"] != nil {
		repoDescription = utils.GenerateFinalComponent((*req.Store)["req:new:repository:description"].(string), req.ExternalData, []string{})
	}
	private := true
	if (*req.Store)["req:new:repository:private"] != nil {
		private = (*req.Store)["req:new:repository:private"].(string) == "true"
	}

	body := make(map[string]interface{})
	body["name"] = repoName
	body["description"] = repoDescription
	body["private"] = private

	ebody, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["CreateNewRepositoryEndpoint"].Call([]interface{}{
		req.Authorization,
		string(ebody),
	})

	return shared.AreaResponse{
		Error: errr,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForGithubReactionCreateNewRepository() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_repository",
		Description: "Create a new repository",
		RequestStore: map[string]static.StoreElement{
			"req:new:repository:name": {
				Description: "The name of the repository to create",
				Type:        "string",
				Required:    true,
			},
			"req:new:repository:description": {
				Description: "The description of the repository to create",
				Type:        "string",
				Required:    false,
			},
			"req:new:repository:private": {
				Description: "Whether the repository should be private or not",
				Type:        "select",
				Required:    false,
				Values: []string{
					"true",
					"false",
				},
			},
		},
		Method: createNewRepository,
	}
}
