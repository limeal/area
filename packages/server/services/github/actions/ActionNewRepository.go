package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"regexp"
	"time"
)

// `GetAllRepositoriesResponse` is a struct with a field `Repositories` of type `[]common.Repository`.
//
// The `[]` indicates that `Repositories` is a slice of `common.Repository`s.
//
// The `common.Repository` type is defined in the `common` package.
//
// The `common` package is imported in the `import` section at the top of the file.
//
// The `common` package is imported from the `github.com/jamesjoshuahill/secret/common` repository
// @property {[]common.Repository} Repositories - An array of Repository objects.
type GetAllRepositoriesResponse struct {
	Repositories []common.Repository `json:"repositories"`
}

// It checks if there's a new repository on GitHub
func hasANewRepository(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)
	query["sort"] = "created"
	query["direction"] = "desc"
	query["type"] = "all"
	query["per_page"] = "1"

	if (*req.Store)["req:repository:type"] != nil {
		query["type"] = (*req.Store)["req:repository:type"].(string)
	}

	if (*req.Store)["ctx:current:time"] == nil {
		(*req.Store)["ctx:current:time"] = time.Now()
	}

	var err error
	var encode []byte
	if (*req.Store)["req:repository:owner"] != nil {
		encode, _, err = req.Service.Endpoints["GetAllRepositoriesFromUserEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			(*req.Store)["req:repository:owner"],
			query,
		})
	} else {
		encode, _, err = req.Service.Endpoints["GetAllRepositoriesEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			query,
		})
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	repositories := GetAllRepositoriesResponse{}
	if err := json.Unmarshal(encode, &repositories.Repositories); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbRepositories := len(repositories.Repositories)
	ok, err := utils.IsLatestByDate(req.Store, nbRepositories, func() interface{} {
		return repositories.Repositories[0].ID
	}, func() string {
		return repositories.Repositories[0].CreatedAt
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["req:repository:regex"] != nil {
		match, err := regexp.MatchString((*req.Store)["req:repository:regex"].(string), repositories.Repositories[0].Name)
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
		if !match {
			(*req.Store)["ctx:last:total"] = nil
			(*req.Store)["ctx:last:elem"] = nil
			return shared.AreaResponse{Success: false}
		}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"github:repository:name":        repositories.Repositories[0].Name,
			"github:repository:owner":       repositories.Repositories[0].Owner.Login,
			"github:repository:description": repositories.Repositories[0].Description,
			"github:repository:url":         repositories.Repositories[0].HTMLURL,
			"github:repository:private":     repositories.Repositories[0].Private,
		},
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForGithubActionNewRepository() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_repository",
		Description: "Triggered when a new repository is created",
		RequestStore: map[string]static.StoreElement{
			"req:repository:regex": {
				Priority:    1,
				Description: "The regex to match the repository name",
				Required:    false,
			},
			"req:repository:owner": {
				Description: "The owner of the repository",
				Required:    false,
			},
			"req:repository:type": {
				Priority:    2,
				Type:        "select",
				Description: "The type of the repository (private | public | member | owner | all) (default: all)",
				Required:    false,
				Values: []string{
					"private",
					"public",
					"member",
					"owner",
					"all",
				},
			},
		},
		Components: []string{
			"github:repository:name",
			"github:repository:owner",
			"github:repository:description",
			"github:repository:url",
			"github:repository:private",
		},
		Method: hasANewRepository,
	}
}
