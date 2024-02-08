package github

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
)

// It returns a map of endpoints for the Github API
func GithubEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Validators
		"GetUserByLoginEndpoint": {
			BaseURL:        "https://api.github.com/users/${login}",
			Params:         GetUserByLoginEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetRepositoryFromNameEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}",
			Params:         GetRepositoryFromNameEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetBranchFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/branches/${branch}",
			Params:         GetBranchFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetReleaseByTagNameEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/releases/tags/${tag}",
			Params:         GetReleaseByTagNameEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetCommitFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/commits/${ref}",
			Params:         GetCommitFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"GetAllRepositoriesEndpoint": {
			BaseURL:        "https://api.github.com/user/repos",
			Params:         GetAllRepositoriesEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllRepositoriesFromUserEndpoint": {
			BaseURL:        "https://api.github.com/users/${login}/repos",
			Params:         GetAllRepositoriesFromUserEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllCollaboratorsFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/collaborators",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllCommitFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/commits",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllReleaseFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/releases",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllPullRequestFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/pulls",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllIssueFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/issues",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetAllBranchFromRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/branches",
			Params:         GetAllFromRepositoryEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"CreateNewGistEndpoint": {
			BaseURL:        "https://api.github.com/gists",
			Params:         CreateGistEndpointParams,
			ExpectedStatus: []int{201},
		},
		"CreateNewIssueEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/issues",
			Params:         CreateInRepositoryEndpointParams,
			ExpectedStatus: []int{201},
		},
		"CreateNewPullRequestEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/pulls",
			Params:         CreateInRepositoryEndpointParams,
			ExpectedStatus: []int{201},
		},
		"AddCollaboratorToRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/collaborators/${collaborator}",
			Params:         AddCollaboratorToRepositoryEndpointParams,
			ExpectedStatus: []int{201, 204},
		},
		"CreateNewReleaseEndpoint": {
			BaseURL:        "https://api.github.com/repos/${owner}/${repo}/releases",
			Params:         CreateInRepositoryEndpointParams,
			ExpectedStatus: []int{201},
		},
		"CreateNewRepositoryEndpoint": {
			BaseURL:        "https://api.github.com/user/repos",
			Params:         CreateNewRepositoryEndpointParams,
			ExpectedStatus: []int{201},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetAllRepositoriesFromUserEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"login": params[1].(string),
		},
		QueryParams: params[2].(map[string]string),
	}
}

// It takes two parameters, the first one is a pointer to a `models.Authorization` struct and the
// second one is a string. It returns a pointer to a `utils.RequestParams` struct
func GetUserByLoginEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"login": params[1].(string),
		},
	}
}

// `GetReleaseByTagNameEndpointParams` is a function that returns a `*utils.RequestParams` object that
// contains the parameters needed to make a request to the `GetReleaseByTagName` endpoint
func GetReleaseByTagNameEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner": params[1].(string),
			"repo":  params[2].(string),
			"tag":   params[3].(string),
		},
	}
}

// It takes in two parameters, the first one is a struct of type `models.Authorization` and the second
// one is a map of type `map[string]string`. It returns a struct of type `utils.RequestParams`
func GetAllRepositoriesEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		QueryParams: params[1].(map[string]string),
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetCommitFromRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner": params[1].(string),
			"repo":  params[2].(string),
			"ref":   params[3].(string),
		},
	}
}

// > This function takes in a slice of interfaces and returns a pointer to a RequestParams struct
func GetRepositoryFromNameEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner": params[1].(string),
			"repo":  params[2].(string),
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetAllFromRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner": params[1].(string),
			"repo":  params[2].(string),
		},
		QueryParams: params[3].(map[string]string),
	}
}

// It takes two parameters, the first one is a pointer to a `models.Authorization` struct and the
// second one is a string. It returns a pointer to a `utils.RequestParams` struct
func CreateGistEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		Body: params[1].(string),
	}
}

// `AddCollaboratorToRepositoryEndpointParams` is a function that takes in a slice of `interface{}` and
// returns a `*utils.RequestParams`
//
// The `interface{}` is a slice of `interface{}` because we want to be able to pass in any number of
// arguments to the function
func AddCollaboratorToRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner":        params[1].(string),
			"repo":         params[2].(string),
			"collaborator": params[3].(string),
		},
	}
}

// `GetBranchFromRepositoryEndpointParams` takes in a slice of `interface{}` and returns a
// `*utils.RequestParams`
//
// The first parameter is the `Authorization` object, which is a struct that contains the `AccessToken`
// string
func GetBranchFromRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner":  params[1].(string),
			"repo":   params[2].(string),
			"branch": params[3].(string),
		},
	}
}

// It creates a request params object for the endpoint `POST /repos/{owner}/{repo}/hooks`
func CreateInRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		UrlParams: map[string]string{
			"owner": params[1].(string),
			"repo":  params[2].(string),
		},
		Body: params[3].(string),
	}
}

// It takes in two parameters, the first being a pointer to a struct of type `models.Authorization` and
// the second being a string. It returns a pointer to a struct of type `utils.RequestParams`
func CreateNewRepositoryEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization":        "Bearer " + params[0].(*models.Authorization).AccessToken,
			"Accept":               "application/vnd.github.v3+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		Body: params[1].(string),
	}
}
