package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"strconv"
)

// It creates a new release for a given repository
func createNewRelease(req static.AreaRequest) shared.AreaResponse {

	repository := utils.GenerateFinalComponent((*req.Store)["req:repository:name"].(string), req.ExternalData, []string{
		"github:repository:name",
	})

	owner := req.AuthStore["login"].(string) // Default value
	if (*req.Store)["req:repository:owner"] != nil {
		owner = utils.GenerateFinalComponent((*req.Store)["req:repository:owner"].(string), req.ExternalData, []string{
			"github:repository:owner",
		})
	}

	release := &common.CreateRelease{
		TagName: utils.GenerateFinalComponent((*req.Store)["req:unused:release:tag:name"].(string), req.ExternalData, []string{}),
	}

	if (*req.Store)["req:release:name"] != nil {
		release.Name = utils.GenerateFinalComponent((*req.Store)["req:release:name"].(string), req.ExternalData, []string{})
	}

	if (*req.Store)["req:release:body"] != nil {
		release.Body = utils.GenerateFinalComponent((*req.Store)["req:release:body"].(string), req.ExternalData, []string{})
	}

	if (*req.Store)["req:release:commitish:value"] != nil {
		release.TargetCommitish = (*req.Store)["req:release:commitish:value"].(string)
	}

	if (*req.Store)["req:release:prerelease"] != nil {
		var err error
		release.Prerelease, err = strconv.ParseBool((*req.Store)["req:release:prerelease"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["req:is:draft"] != nil {
		var err error
		release.Draft, err = strconv.ParseBool((*req.Store)["req:is:draft"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	str, errC := json.Marshal(release)
	if errC != nil {
		return shared.AreaResponse{Error: errC}
	}

	_, _, err := req.Service.Endpoints["CreateNewReleaseEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		owner,
		repository,
		string(str),
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It creates a new release on a github repository
func DescriptorForGithubReactionCreateNewRelease() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_release",
		Description: "Create a new release",
		RequestStore: map[string]static.StoreElement{
			"req:repository:name": {
				Priority:    1,
				Type:        "select_uri",
				Description: "The name of the repository",
				Values:      []string{"/${req:repository:owner}/repos"},
				Required:    true,
			},
			"req:unused:release:tag:name": {
				Priority:    2,
				Type:        "string",
				Description: "The name of the tag (must be unique)",
				Required:    true,
			},
			// Optional
			"req:repository:owner": {
				Type:        "string",
				Description: "The owner of the repository (default: you)",
				Required:    false,
			},
			"req:release:name": {
				Priority:    3,
				Type:        "string",
				Description: "The name of the release (default: tag name)",
				Required:    false,
			},
			"req:release:body": {
				Priority:    4,
				Type:        "string",
				Description: "The body of the release (default: empty)",
				Required:    false,
			},
			"req:release:commitish:type": {
				Priority:    5,
				Type:        "select",
				Description: "Type of the commitish (branch or commit)",
				Required:    false,
				Values:      []string{"branch", "commit"},
			},
			"req:release:commitish:value": {
				Priority:    6,
				Type:        "select_uri",
				Description: "The commitish of the release (commit sha or branch name)",
				Required:    false,
				NeedFields:  []string{"req:release:commitish:type"},
				Values:      []string{"/${req:repository:owner}/${req:repository:name}/${req:release:commitish:type}s"},
			},
			"req:release:prerelease": {
				Priority:    7,
				Type:        "select",
				Description: "Is it a pre-release ? (default: false)",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:is:draft": {
				Priority:    7,
				Type:        "select",
				Description: "Is it a draft ?  (default: false)",
				Required:    false,
				Values:      []string{"true", "false"},
			},
		},
		Method: createNewRelease,
	}
}
