package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"regexp"
)

// TODO: Minor regex branch not working

// `GetAllBranchFromRepositoryResponse` is a struct with a field `Branches` of type `[]common.Branch`.
// @property {[]common.Branch} Branches - An array of branches.
type GetAllBranchFromRepositoryResponse struct {
	Branches []common.Branch `json:"branches"`
}

// It checks if a new branch has been created in a repository
func hasANewBranch(req static.AreaRequest) shared.AreaResponse {

	userLogin := req.AuthStore["login"]

	if (*req.Store)["req:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"].(string)
	}

	query := make(map[string]string)
	query["per_page"] = "1"

	if (*req.Store)["req:only:protected"] != nil {
		query["protected"] = (*req.Store)["req:only:protected"].(string)
	}

	encode, httpResp, err := req.Service.Endpoints["GetAllBranchFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	branches := GetAllBranchFromRepositoryResponse{}
	if err := json.Unmarshal(encode, &branches.Branches); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbBranches, errP := common.GetLastPage(httpResp.Header.Get("Link"))
	if errP != nil {
		return shared.AreaResponse{Error: errP}
	}
	ok, errL := utils.IsLatestBasic(req.Store, nbBranches)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:branch:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:branch:regex"].(string), branches.Branches[0].Name)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New branch found (Repo: "+(*req.Store)["req:repository:name"].(string)+") (Name: "+branches.Branches[0].Name+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:repository:name":  (*req.Store)["req:repository:name"],
			"github:repository:owner": userLogin,
			"github:branch:name":      branches.Branches[0].Name,
			"github:commit:sha":       branches.Branches[0].Commit.Sha,
			"github:commit:url":       branches.Branches[0].Commit.Url,
			"github:branch:protected": branches.Branches[0].Protected,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_branch" and the method that
// should be called to check if the service area is triggered is the hasANewBranch function
func DescriptorForGithubActionAnyNewBranch() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_branch",
		Description: "When a new branch is created",
		Method:      hasANewBranch,
		RequestStore: map[string]static.StoreElement{
			"req:repository:name": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The name of the repository",
				Values:      []string{"/${req:repository:owner}/repos"},
				Required:    true,
			},
			"req:only:protected": {
				Priority:    2,
				Type:        "select",
				Description: "Only trigger if the branch is protected",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository",
				Required:    false,
			},
			"req:branch:regex": {
				Priority:    2,
				Type:        "string",
				Description: "The regex to match the branch name",
				Required:    false,
			},
		},
		Components: []string{
			"github:repository:name",
			"github:repository:owner",
			"github:branch:name",
			"github:commit:sha",
			"github:commit:url",
			"github:branch:protected",
		},
	}
}
