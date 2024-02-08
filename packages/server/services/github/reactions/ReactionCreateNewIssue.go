package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"strconv"
	"strings"
)

// It creates a new issue in a repository
func createNewIssue(req static.AreaRequest) shared.AreaResponse {
	// Find the user the first time

	repository := utils.GenerateFinalComponent((*req.Store)["req:repository:name"].(string), req.ExternalData, []string{
		"github:repository:name",
	})
	owner := req.AuthStore["login"]
	if (*req.Store)["req:repository:owner"] != nil {
		owner = utils.GenerateFinalComponent((*req.Store)["req:repository:owner"].(string), req.ExternalData, []string{
			"github:repository:owner",
		})
	}
	title := utils.GenerateFinalComponent((*req.Store)["req:issue:title"].(string), req.ExternalData, []string{})
	body := utils.GenerateFinalComponent((*req.Store)["req:issue:body"].(string), req.ExternalData, []string{})

	issue := &common.CreateIssue{
		Title: title,
		Body:  body,
	}

	if (*req.Store)["req:issue:assignees"] != nil {
		issue.Assignees = strings.Split((*req.Store)["req:issue:assignees"].(string), ",")
	}

	if (*req.Store)["req:issue:labels"] != nil {
		issue.Labels = strings.Split((*req.Store)["req:issue:labels"].(string), ",")
	}

	if (*req.Store)["req:issue:milestone"] != nil {
		var err error
		issue.Milestone, err = strconv.Atoi((*req.Store)["req:issue:milestone"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	str, err := json.Marshal(issue)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, errr := req.Service.Endpoints["CreateNewIssueEndpoint"].Call([]interface{}{
		req.Authorization,
		owner,
		repository,
		string(str),
	})
	if errr != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForGithubReactionCreateNewIssue() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_issue",
		Description: "Create a new issue",
		Method:      createNewIssue,
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
			"req:issue:title": {
				Priority:    2,
				Description: "The title of the issue",
				Required:    true,
			},
			"req:issue:body": {
				Priority:    3,
				Description: "The body of the issue",
				Required:    true,
			},
			"req:issue:assignees": {
				Priority:    4,
				Type:        "long_string",
				Description: "The assignees of the issue",
				Required:    false,
			},
			"req:issue:labels": {
				Priority:    4,
				Type:        "long_string",
				Description: "The labels of the issue",
				Required:    false,
			},
			"req:issue:milestone": {
				Priority:    4,
				Description: "The milestone of the issue",
				Required:    false,
			},
		},
	}
}
