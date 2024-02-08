package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"strconv"
)

// `GetAllPullRequestFromRepositoryResponse` is a struct with a field `PullRequests` of type
// `[]common.PullRequest`.
// @property {[]common.PullRequest} PullRequests - An array of PullRequest objects.
type GetAllPullRequestFromRepositoryResponse struct {
	PullRequests []common.PullRequest `json:"pull_requests"`
}

// It checks if there is a new pull request on a repository
func hasANewPullRequest(req static.AreaRequest) shared.AreaResponse {
	// Find the user the first time
	userLogin := req.AuthStore["login"]

	if (*req.Store)["ctx:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"]
	}

	query := make(map[string]string)
	query["sort"] = "created"
	query["direction"] = "desc"
	query["per_page"] = "1"

	encode, _, err := req.Service.Endpoints["GetAllPullRequestFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	pullRequests := GetAllPullRequestFromRepositoryResponse{}
	if err := json.Unmarshal(encode, &pullRequests.PullRequests); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbPullRequests := len(pullRequests.PullRequests)
	ok, err := utils.IsLatestByID(req.Store, nbPullRequests, func() int {
		return pullRequests.PullRequests[0].ID
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	req.Logger.WriteInfo("[Action] New pull request found (Repo: "+(*req.Store)["req:repository:name"].(string)+") (Number: "+strconv.Itoa(pullRequests.PullRequests[0].Number)+")", true)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:pull:number":      pullRequests.PullRequests[0].Number,
			"github:pull:title":       pullRequests.PullRequests[0].Title,
			"github:pull:body":        pullRequests.PullRequests[0].Body,
			"github:pull:state":       pullRequests.PullRequests[0].State,
			"github:pull:html":        pullRequests.PullRequests[0].HTMLURL,
			"github:author:login":     pullRequests.PullRequests[0].User.Login,
			"github:repository:name":  pullRequests.PullRequests[0].Head.Repo.Name,
			"github:repository:owner": pullRequests.PullRequests[0].Head.Repo.Owner.Login,
			"github:branch:name":      pullRequests.PullRequests[0].Head.Ref,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_pull_request" and the method
// that will be called to check if the service area is triggered is the function "hasANewPullRequest"
func DescriptorForGithubActionAnyNewPullRequest() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_pull_request",
		Description: "When a new pull request is created",
		Method:      hasANewPullRequest,
		RequestStore: map[string]static.StoreElement{
			"req:repository:name": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The name of the repository",
				Required:    true,
				Values:      []string{"/${req:repository:owner}/repos"},
			},
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository",
				Required:    false,
			},
		},
		Components: []string{
			"github:pull:number",
			"github:pull:title",
			"github:pull:body",
			"github:pull:state",
			"github:pull:html",
			"github:author:login",
			"github:repository:name",
			"github:repository:owner",
			"github:branch:name",
		},
	}
}
