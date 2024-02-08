package github

import (
	"area-server/classes/static"
	"area-server/services/github/common"
	"area-server/utils"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of static.ServiceRoute, which is a struct that contains the endpoint, the
// handler, the method and a boolean that indicates if the route needs authentication
func GithubRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/:owner/repos",
			Handler:  GetUserRepositoriesRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/:owner/:repo/branchs",
			Handler:  GetBranchesFromRepositoryRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/:owner/:repo/commits",
			Handler:  GetCommitsFromRepositoryRoute,
			Method:   "GET",
			NeedAuth: true,
		},
	}
}

// ---------------------- Github Routes ----------------------

// `GetUserRepositoriesResponse` is a struct with a single field, `Items`, which is a slice of
// `common.Repository`s.
// @property {[]common.Repository} Items - An array of repositories that the user has access to.
type GetUserRepositoriesResponse struct {
	Items []common.Repository `json:"items"`
}

// It gets all the repositories from a user
func GetUserRepositoriesRoute(c *fiber.Ctx) error {

	auth0, errO := utils.VerifyRoute(c, "github")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()
	owner := c.Params("owner")
	typev := c.Query("type", "all")
	sort := c.Query("sort", "created")

	var err error
	var encode []byte
	if owner == "default" {
		encode, _, err = service.Endpoints["GetAllRepositoriesEndpoint"].CallEncode([]interface{}{
			auth0,
			map[string]string{
				"type": typev,
				"sort": sort,
			},
		})
	} else {
		encode, _, err = service.Endpoints["GetAllRepositoriesFromUserEndpoint"].CallEncode([]interface{}{
			auth0,
			owner,
			map[string]string{
				"type": typev,
				"sort": sort,
			},
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	repositories := GetUserRepositoriesResponse{}
	err = json.Unmarshal(encode, &repositories.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   repositories.Items,
			"fields": []string{"name", "name"},
		},
	})
}

// `GetBranchesFromRepositoryResponse` is a struct with a single field, `Items`, which is a slice of
// `common.Branch`s.
// @property {[]common.Branch} Items - An array of branches.
type GetBranchesFromRepositoryResponse struct {
	Items []common.Branch `json:"items"`
}

// It gets all the branches from a repository
func GetBranchesFromRepositoryRoute(c *fiber.Ctx) error {

	auth0, errO := utils.VerifyRoute(c, "github")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	authStore := make(map[string]interface{})
	if err := json.Unmarshal([]byte(auth0.Other.String()), &authStore); err != nil {
		return errors.New("Area: Auth Store is not valid !")
	}

	service := Descriptor()
	owner := c.Params("owner", "")
	repo := c.Params("repo")

	if owner == "default" {
		owner = authStore["login"].(string)
	}

	query := make(map[string]string)
	query["per_page"] = "15"

	encode, _, err := service.Endpoints["GetAllBranchFromRepositoryEndpoint"].CallEncode([]interface{}{
		auth0,
		owner,
		repo,
		query,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	branches := GetBranchesFromRepositoryResponse{}
	err = json.Unmarshal(encode, &branches.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   branches.Items,
			"fields": []string{"name", "name"},
		},
	})
}

// `GetCommitsFromRepositoryResponse` is a struct with a single field, `Items`, which is a slice of
// `common.Commit`s.
// @property {[]common.Commit} Items - An array of commits.
type GetCommitsFromRepositoryResponse struct {
	Items []common.Commit `json:"items"`
}

// It gets all commits from a repository
func GetCommitsFromRepositoryRoute(c *fiber.Ctx) error {

	auth0, errO := utils.VerifyRoute(c, "github")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	authStore := make(map[string]interface{})
	if err := json.Unmarshal([]byte(auth0.Other.String()), &authStore); err != nil {
		return errors.New("Area: Auth Store is not valid !")
	}

	service := Descriptor()
	owner := c.Params("owner")
	repo := c.Params("repo")

	if owner == "default" {
		owner = authStore["login"].(string)
	}

	query := make(map[string]string)
	query["per_page"] = "15"

	encode, _, err := service.Endpoints["GetAllCommitFromRepositoryEndpoint"].CallEncode([]interface{}{
		auth0,
		owner,
		repo,
		query,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	commits := GetCommitsFromRepositoryResponse{}
	err = json.Unmarshal(encode, &commits.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   commits.Items,
			"fields": []string{"sha", "sha"},
		},
	})
}
