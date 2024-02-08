package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"encoding/json"
	"regexp"
	"time"
)

// `GetAllCommitFromRepositoryResponse` is a struct with a field `Commits` of type `[]common.Commit`.
// @property {[]common.Commit} Commits - An array of commits.
type GetAllCommitFromRepositoryResponse struct {
	Commits []common.Commit `json:"commits"`
}

// It checks if there is a new commit on a repository
func hasANewCommit(req static.AreaRequest) shared.AreaResponse {
	userLogin := req.AuthStore["login"]

	if (*req.Store)["req:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"].(string)
	}

	if (*req.Store)["ctx:current:time"] == nil {
		(*req.Store)["ctx:current:time"] = time.Now().UTC().Format("2006-01-02T15:04:05Z0700")
	}

	query := make(map[string]string)
	query["per_page"] = "1"
	query["since"] = (*req.Store)["ctx:current:time"].(string)

	if (*req.Store)["req:branch:name"] != nil {
		query["sha"] = (*req.Store)["req:branch:name"].(string)
	}

	if (*req.Store)["req:commit:author"] != nil {
		query["author"] = (*req.Store)["req:commit:author"].(string)
	}

	encode, _, err := req.Service.Endpoints["GetAllCommitFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	commits := GetAllCommitFromRepositoryResponse{}
	if err := json.Unmarshal(encode, &commits.Commits); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbCommits := len(commits.Commits)
	if nbCommits == 0 {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:commit:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:commit:regex"].(string), commits.Commits[0].ICommit.Message)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			return shared.AreaResponse{Success: false}
		}
	}

	req.Logger.WriteInfo("[Action] New commit found (Repo: "+(*req.Store)["req:repository:name"].(string)+") (Sha: "+commits.Commits[0].Sha+")", false)
	(*req.Store)["ctx:current:time"] = time.Now().UTC().Format("2006-01-02T15:04:05Z0700")
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:repository:name":  (*req.Store)["req:repository:name"],
			"github:repository:owner": userLogin,
			"github:commit:sha":       commits.Commits[0].Sha,
			"github:commit:msg":       commits.Commits[0].ICommit.Message,
			"github:author:login":     commits.Commits[0].ICommit.Author.Name,
			"github:author:email":     commits.Commits[0].ICommit.Author.Email,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_commit" that is triggered when
// a new commit is pushed
func DescriptorForGithubActionAnyNewCommit() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_commit",
		Description: "When a new commit is pushed",
		Method:      hasANewCommit,
		RequestStore: map[string]static.StoreElement{
			// Required
			"req:repository:name": {
				Priority:    1,
				Type:        "select_uri",
				Description: "Name of the repository to watch",
				Values:      []string{"/${req:repository:owner}/repos"},
				Required:    true,
			},
			// Optional
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository",
				Required:    false,
			},
			"req:branch:name": {
				Priority:    2,
				Type:        "select_uri",
				Description: "The branch to watch",
				NeedFields:  []string{"req:repository:name"},
				Required:    false,
				Values:      []string{"/${req:repository:owner}/${req:repository:name}/branchs"},
			},
			"req:commit:author": {
				Priority:    2,
				Type:        "string",
				Description: "The author of the commit",
				Required:    false,
			},
			"req:commit:regex": {
				Priority:    2,
				Type:        "string",
				Description: "The regex to match the commit message",
				Required:    false,
			},
		},
		Components: []string{
			"github:repository:name",
			"github:repository:owner",
			"github:commit:sha",
			"github:commit:msg",
			"github:author:login",
			"github:author:email",
		},
	}
}
