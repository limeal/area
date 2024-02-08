package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"errors"
)

// It adds a collaborator to a repository
func addCollaborator(req static.AreaRequest) shared.AreaResponse {
	repository := utils.GenerateFinalComponent((*req.Store)["req:repository:name"].(string), req.ExternalData, []string{
		"github:repository:name",
	})
	owner := req.AuthStore["login"]
	if (*req.Store)["req:repository:owner"] != nil {
		owner = utils.GenerateFinalComponent((*req.Store)["req:repository:owner"].(string), req.ExternalData, []string{
			"github:repository:owner",
		})
	}

	_, httpResp, err := req.Service.Endpoints["AddCollaboratorToRepositoryEndpoint"].Call([]interface{}{
		req.Authorization,
		owner,
		repository,
		(*req.Store)["req:collaborator:login"].(string),
	})

	if httpResp != nil && httpResp.StatusCode == 422 {
		return shared.AreaResponse{Error: errors.New("Cannot add yourself as a collaborator")}
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if httpResp.StatusCode == 204 {
		req.Logger.WriteInfo("[Reaction] Collaborator already added or it is useless to add it (Collaborator Login:"+(*req.Store)["req:collaborator:login"].(string)+" )", false)
	} else {
		req.Logger.WriteInfo("[Reaction] Collaborator added (Collaborator Login:"+(*req.Store)["req:collaborator:login"].(string)+" )", false)
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForGithubReactionAddCollaborator() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_collaborator",
		Description: "When a new collaborator is added to a repository",
		RequestStore: map[string]static.StoreElement{
			"req:repository:name": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The name of the repository",
				Values:      []string{"/${req:repository:owner}/repos"},
				Required:    true,
			},
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository",
				Required:    false,
			},
			"req:collaborator:login": {
				Priority:    2,
				Type:        "string",
				Description: "The login of the collaborator",
				Required:    true,
			},
			"req:collaborator:permission": {
				Priority:    3,
				Type:        "select",
				Description: "The permission of the collaborator",
				Required:    false,
				Values:      []string{"pull", "push", "admin", "maintain", "triage"},
			},
		},
		Method: addCollaborator,
	}
}
