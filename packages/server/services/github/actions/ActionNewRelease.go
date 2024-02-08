package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
)

// `GetAllReleaseFromRepositoryResponse` is a struct with a field `Releases` of type
// `[]common.Release`.
// @property {[]common.Release} Releases - An array of releases.
type GetAllReleaseFromRepositoryResponse struct {
	Releases []common.Release `json:"releases"`
}

// It checks if there is a new release for a given repository
func hasANewRelease(req static.AreaRequest) shared.AreaResponse {
	// Find the user the first time
	userLogin := req.AuthStore["login"]

	if (*req.Store)["ctx:repository:owner"] != nil {
		userLogin = (*req.Store)["req:repository:owner"]
	}

	query := make(map[string]string)
	query["per_page"] = "1"

	encode, httpResp, err := req.Service.Endpoints["GetAllReleaseFromRepositoryEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userLogin,
		(*req.Store)["req:repository:name"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	releases := GetAllReleaseFromRepositoryResponse{}
	if err := json.Unmarshal(encode, &releases.Releases); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbReleases, errP := common.GetLastPage(httpResp.Header.Get("Link"))
	if errP != nil {
		return shared.AreaResponse{Error: errP}
	}
	ok, errL := utils.IsLatestBasic(req.Store, nbReleases)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	req.Logger.WriteInfo("[Action] New release found (Repo: "+(*req.Store)["req:repository:name"].(string)+") (Name: "+releases.Releases[0].Name+")", false)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:repository:name":  (*req.Store)["req:repository:name"],
			"github:repository:owner": userLogin,
			"github:release:name":     releases.Releases[0].Name,
			"github:release:id":       releases.Releases[0].ID,
			"github:release:tag":      releases.Releases[0].TagName,
			"github:release:body":     releases.Releases[0].Body,
			"github:release:html":     releases.Releases[0].HTMLURL,
			"github:release:tar":      releases.Releases[0].TarBallURL,
		},
	}
}

// It returns a static.ServiceArea that describes the service area "new_release" and the method that
// will be called to check if there is a new release is the function hasANewRelease
func DescriptorForGithubActionAnyNewRelease() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_release",
		Description: "When a new release is pushed",
		Method:      hasANewRelease,
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
			"github:repository:name",
			"github:repository:owner",
			"github:release:name",
			"github:release:id",
			"github:release:tag",
			"github:release:body",
			"github:release:html",
			"github:release:tar",
		},
	}
}
