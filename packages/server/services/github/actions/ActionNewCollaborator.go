package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
)

// `GetAllCollaboratorsResponse` is a struct with a field `Collaborators` of type
// `[]common.Collaborator`.
//
// The `[]` indicates that `Collaborators` is a slice of `common.Collaborator`s.
//
// The `common.Collaborator` type is defined in the `common` package.
//
// The `common` package is imported in the `import` section of the file.
//
// The `import` section is at the top of the file.
//
// The `import` section is a list of
// @property {[]common.Collaborator} Collaborators - An array of collaborators.
type GetAllCollaboratorsResponse struct {
	Collaborators []common.Collaborator `json:"collaborators"`
}

// It checks if a new collaborator has been added to a repository
func hasANewCollaborator(req static.AreaRequest) shared.AreaResponse {

	userLogin := req.AuthStore["login"]

	if (*req.Store)["req:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"].(string)
	}

	query := make(map[string]string)
	query["per_page"] = "1"

	if (*req.Store)["req:collaborator:permission"] != nil {
		query["permission"] = (*req.Store)["req:collaborator:permission"].(string)
	}

	if (*req.Store)["req:collaborator:affiliation"] != nil {
		query["affiliation"] = (*req.Store)["req:collaborator:affiliation"].(string)
	}

	encode, httpResp, err := req.Service.Endpoints["GetAllCollaboratorsFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	collaborators := GetAllCollaboratorsResponse{}
	err = json.Unmarshal(encode, &collaborators.Collaborators)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbCollaborators, errP := common.GetLastPage(httpResp.Header.Get("Link"))
	if errP != nil {
		return shared.AreaResponse{Error: errP}
	}
	ok, errL := utils.IsLatestBasic(req.Store, nbCollaborators)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:collaborator:login"] != nil {
		if (*req.Store)["req:collaborator:login"].(string) != collaborators.Collaborators[0].Login {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New collaborator in repository (Repository Name: "+(*req.Store)["req:repository:name"].(string)+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:repository:name":    (*req.Store)["req:repository:name"],
			"github:repository:owner":   userLogin,
			"github:collaborator:login": collaborators.Collaborators[0].Login,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_collaborator" and the method
// that will be called to check if the service area is triggered is hasANewCollaborator
func DescriptorForGithubActionNewCollaborator() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_collaborator",
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
				Required:    false,
			},
			"req:collaborator:permission": {
				Priority:    2,
				Type:        "select",
				Description: "The permission of the collaborator (default: all)",
				Required:    false,
				Values:      []string{"admin", "push", "pull", "maintain", "triage"},
			},
			"req:collaborator:affiliation": {
				Priority:    2,
				Type:        "select",
				Description: "The affiliation of the collaborator (default: all)",
				Required:    false,
				Values:      []string{"outside", "direct", "all"},
			},
		},
		Method: hasANewCollaborator,
		Components: []string{
			"github:repository:name",
			"github:repository:owner",
			"github:collaborator:login",
		},
	}
}
