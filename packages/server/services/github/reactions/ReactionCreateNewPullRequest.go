package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"strconv"
)

// It creates a new pull request
func createNewPullRequest(req static.AreaRequest) shared.AreaResponse {

	repoName := utils.GenerateFinalComponent((*req.Store)["req:repository:name"].(string), req.ExternalData, []string{
		"github:repository:name",
	})
	repoOwner := req.AuthStore["login"]
	if (*req.Store)["req:repository:owner"] != nil {
		repoOwner = utils.GenerateFinalComponent((*req.Store)["req:repository:owner"].(string), req.ExternalData, []string{
			"github:repository:owner",
		})
	}

	title := utils.GenerateFinalComponent((*req.Store)["req:pull:title"].(string), req.ExternalData, []string{})
	body := ""
	if (*req.Store)["req:pull:body"] != nil {
		body = utils.GenerateFinalComponent((*req.Store)["req:pull:body"].(string), req.ExternalData, []string{})
	}
	branchSource := utils.GenerateFinalComponent((*req.Store)["req:branch:source"].(string), req.ExternalData, []string{
		"github:branch:name",
	})
	branchTarget := utils.GenerateFinalComponent((*req.Store)["req:branch:target"].(string), req.ExternalData, []string{
		"github:branch:name",
	})

	headRepo := repoName
	if (*req.Store)["req:head:repository:name"] != nil {
		headRepo = utils.GenerateFinalComponent((*req.Store)["req:head:repository:name"].(string), req.ExternalData, []string{
			"github:repository:name",
		})
	}

	draftMode := false
	if (*req.Store)["req:is:draft"] != nil && (*req.Store)["req:is:draft"].(string) == "true" {
		draftMode = true
	}

	rbody := make(map[string]interface{})
	rbody["head"] = branchSource
	rbody["head_repo"] = headRepo
	rbody["base"] = branchTarget
	rbody["body"] = body
	rbody["draft"] = draftMode

	if (*req.Store)["req:from:issue"] != nil && (*req.Store)["req:from:issue"].(string) == "true" {
		var err error
		rbody["issue"], err = strconv.Atoi(title)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	} else {
		rbody["title"] = title
	}

	ebody, err := json.Marshal(rbody)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["CreateNewPullRequestEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		repoOwner,
		repoName,
		string(ebody),
	})

	return shared.AreaResponse{
		Error: errr,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForGithubReactionCreateNewPullRequest() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_pull_request",
		Description: "Create a new pull request",
		RequestStore: map[string]static.StoreElement{
			"req:repository:name": {
				Priority:    3,
				Type:        "select_uri",
				Description: "The name of the repository",
				Values:      []string{"/${req:repository:owner}/repos"},
				Required:    true,
			},
			"req:pull:title": {
				Priority:    1,
				Type:        "string",
				Description: "The title of the pull request",
				Required:    true,
			},
			"req:from:issue": {
				Priority:    2,
				Type:        "select",
				Description: "If true, the pull request will be created from an issue with name in req:pull:title",
				Required:    true,
				Values:      []string{"true", "false"},
			},
			"req:branch:source": {
				Priority:    4,
				Type:        "select_uri",
				Description: "The source branch",
				Required:    true,
				Values:      []string{"/${req:repository:owner}/${req:repository:name}/branchs"},
			},
			"req:branch:target": {
				Priority:    6,
				Type:        "select_uri",
				Description: "The target branch (must be different from the source branch)",
				Required:    true,
				Values:      []string{"/${req:repository:owner}/${req:repository:name}/branchs"},
			},
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository",
				Required:    false,
			},
			"req:head:repository:name": {
				Priority:    5,
				Type:        "select_uri",
				Description: "The name of the repository where the source branch is located. If empty, the source branch is located in the same repository as the target branch.",
				Required:    false,
				Values:      []string{"/${req:repository:owner}/repos"},
			},
			"req:pull:body": {
				Priority:    7,
				Type:        "long_string",
				Description: "The body of the pull request",
				Required:    false,
			},
			"req:is:draft": {
				Priority:    7,
				Type:        "select",
				Description: "If true, the pull request will be created as a draft",
				Required:    false,
				Values:      []string{"true", "false"},
			},
		},
		Method: createNewPullRequest,
	}
}
