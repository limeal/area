package github

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/common"
	"encoding/json"
)

// It returns a map of validators for the Github service
func GithubValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:repository:owner":         UserValidator,
		"req:repository:name":          RepoNameValidator,
		"req:gist:public":              common.BoolValidator,
		"req:only:protected":           common.BoolValidator,
		"req:repository:type":          RepoTypeValidator,
		"req:collaborator:login":       UserValidator,
		"req:collaborator:permission":  PermissionValidator,
		"req:collaborator:affiliation": AffiliationValidator,
		"req:commit:author":            UserValidator,
		"req:branch:sha":               BranchValidator,
		"req:unused:release:tag:name":  TagValidator,
		"req:release:commitish:type":   CommitishTypeValidator,
		"req:release:commitish:value":  CommitishValueValidator,
		"req:release:prerelease":       common.BoolValidator,
		"req:is:draft":                 common.BoolValidator,
		"req:issue:milestone":          common.IntValidator,
	}
}

// ----------------------- Validators -----------------------

// It takes an authorization, a service, a value, and a store, and returns a boolean
func UserValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	_, _, err := service.Endpoints["GetUserByLoginEndpoint"].CallEncode([]interface{}{auth, value.(string)})
	if err != nil {
		return false
	}

	return true
}

// It takes an authorization, a service, a value, and a store, and returns true if the value is a
// string and the value is a valid repository name for the user
func RepoNameValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	owner := store["req:repository:owner"]

	if owner == nil {
		var other map[string]interface{}
		err := json.Unmarshal(auth.Other, &other)
		if err != nil {
			return false
		}
		owner = other["login"]
	}

	if owner == nil {
		return false
	}

	_, _, err := service.Endpoints["GetRepositoryFromNameEndpoint"].CallEncode([]interface{}{auth, owner, value.(string)})
	if err != nil {
		return false
	}

	return true
}

// If the value is a string, and it's one of the following values: "public", "private", "member",
// "owner", or "all", then it's valid
func RepoTypeValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "public" || value == "private" || value == "member" || value == "owner" || value == "all" {
		return true
	}

	return false
}

// If the value is a string, and the string is one of the valid permissions, then return true
func PermissionValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "admin" || value == "push" || value == "pull" || value == "maintain" || value == "triage" {
		return true
	}

	return false
}

// If the value is "outside", "direct", or "all", then it's valid
func AffiliationValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "outside" || value == "direct" || value == "all" {
		return true
	}

	return false
}

// It checks if the value is a string, and if it is, it checks if the value is a valid branch in the
// repository
func BranchValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:repository:name"] == nil {
		return false
	}

	owner := store["req:repository:owner"]

	if owner == nil {
		var other map[string]interface{}
		err := json.Unmarshal(auth.Other, &other)
		if err != nil {
			return false
		}
		owner = other["login"]
	}

	if owner == nil {
		return false
	}

	_, _, err := service.Endpoints["GetBranchFromRepositoryEndpoint"].CallEncode([]interface{}{auth, owner, store["req:repository:name"], value.(string)})
	if err != nil {
		return false
	}

	return true
}

// If the value is a string, and the repository name is in the store, and the owner is in the store or
// in the authorization's `other` field, and the GitHub API returns an error when we try to get the
// release by tag name, then the value is valid
func TagValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:repository:name"] == nil {
		return false
	}

	owner := store["req:repository:owner"]

	if owner == nil {
		var other map[string]interface{}
		err := json.Unmarshal(auth.Other, &other)
		if err != nil {
			return false
		}
		owner = other["login"]
	}

	if owner == nil {
		return false
	}

	_, _, err := service.Endpoints["GetReleaseByTagNameEndpoint"].CallEncode([]interface{}{
		auth,
		owner,
		store["req:repository:name"],
		value.(string),
	})

	if err != nil {
		return true
	}

	return false
}

// `CommitishTypeValidator` validates that the value of the `commitish_type` parameter is either
// `branch` or `commit`
func CommitishTypeValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:release:commitish:value"] == nil {
		return false
	}

	if value == "branch" || value == "commit" {
		return true
	}

	return false
}

// It checks if the commitish value is valid for the given repository
func CommitishValueValidator(auth *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if store["req:release:commitish:type"] == nil {
		return false
	}

	if store["req:repository:name"] == nil {
		return false
	}

	owner := store["req:repository:owner"]

	if owner == nil {

		var other map[string]interface{}
		err := json.Unmarshal(auth.Other, &other)
		if err != nil {
			return false
		}
		owner = other["login"]
	}

	if owner == nil {
		return false
	}

	if store["req:release:commitish:type"] == "branch" {
		_, _, err := service.Endpoints["GetBranchFromRepositoryEndpoint"].CallEncode([]interface{}{auth, owner, store["req:repository:name"], value.(string)})
		if err != nil {
			return false
		}
	} else {
		_, _, err := service.Endpoints["GetCommitFromRepositoryEndpoint"].CallEncode([]interface{}{auth, owner, store["req:repository:name"], value.(string)})
		if err != nil {
			return false
		}
	}

	return true
}
