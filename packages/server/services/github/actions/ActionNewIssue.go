package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"regexp"
	"strconv"
)

// `GetAllIssueFromRepositoryResponse` is a struct with a single field `Issues` of type
// `[]common.Issue`.
// @property {[]common.Issue} Issues - This is an array of Issue structs.
type GetAllIssueFromRepositoryResponse struct {
	Issues []common.Issue `json:"issues"`
}

// It checks if there's a new issue on a repository
func hasANewIssue(req static.AreaRequest) shared.AreaResponse {
	// Find the user the first time
	userLogin := req.AuthStore["login"]

	if (*req.Store)["req:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"]
	}

	query := make(map[string]string)
	query["sort"] = "created"
	query["direction"] = "desc"
	query["per_page"] = "1"
	query["state"] = "open"
	query["filter"] = "all"

	if (*req.Store)["req:issue:state"] != nil {
		query["state"] = (*req.Store)["req:issue:state"].(string)
	}

	if (*req.Store)["req:issue:filter"] != nil {
		query["filter"] = (*req.Store)["req:issue:filter"].(string)
	}

	encode, _, err := req.Service.Endpoints["GetAllIssueFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	issues := GetAllIssueFromRepositoryResponse{}
	if err := json.Unmarshal(encode, &issues.Issues); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbIssues := len(issues.Issues)
	ok, err := utils.IsLatestByID(req.Store, nbIssues, func() int {
		return issues.Issues[0].ID
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:issue:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:issue:regex"].(string), issues.Issues[0].Title)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New issue found (Repo: "+(*req.Store)["req:repository:name"].(string)+") (Number: "+strconv.Itoa(issues.Issues[0].Number)+")", true)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:issue:title":      issues.Issues[0].Title,
			"github:issue:body":       issues.Issues[0].Body,
			"github:issue:state":      issues.Issues[0].State,
			"github:author:login":     issues.Issues[0].User.Login,
			"github:repository:name":  issues.Issues[0].Repository.Name,
			"github:repository:owner": issues.Issues[0].Repository.Owner.Login,
			"github:issue:number":     issues.Issues[0].Number,
			"github:issue:url":        issues.Issues[0].HTMLURL,
			"github:issue:labels":     issues.Issues[0].Labels,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_issue" and the method to call
// to check if there is a new issue
func DescriptorForGithubActionAnyNewIssue() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_issue",
		Description: "When a new issue is created",
		Method:      hasANewIssue,
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
			"req:issue:state": {
				Priority:    2,
				Type:        "select",
				Description: "State of the issue to check (default: all)",
				Required:    false,
				Values: []string{
					"all",
					"open",
					"closed",
				},
			},
			"req:issue:filter": {
				Priority:    2,
				Type:        "select",
				Description: "Filter of the issue to check",
				Required:    false,
				Values: []string{
					"all",
					"assigned",
					"created",
					"mentioned",
					"subscribed",
					"repos",
				},
			},
			"req:issue:regex": {
				Priority:    2,
				Type:        "string",
				Description: "Regex to match the issue title",
				Required:    false,
			},
		},
		Components: []string{
			"github:issue:title",
			"github:issue:body",
			"github:issue:state",
			"github:author:login",
			"github:repository:name",
			"github:repository:owner",
			"github:issue:number",
			"github:issue:url",
			"github:issue:labels",
		},
	}
}
