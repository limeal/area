package dropbox

import (
	"area-server/classes/static"
	"area-server/services/dropbox/common"
	"area-server/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of `static.ServiceRoute` structs
func DropboxRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/directories",
			Method:   "GET",
			Handler:  ListFoldersHandler,
			NeedAuth: true,
		},
		{
			Endpoint: "/files",
			Method:   "GET",
			Handler:  ListFilesHandler,
			NeedAuth: true,
		},
	}
}

// ------------------ Routes ------------------

// `ListFoldersResponse` is a struct with three fields: `Entries`, `Cursor`, and `HasMore`.
//
// The `Entries` field is an array of `common.Entry`s. The `Cursor` field is a string. The `HasMore`
// field is a boolean.
// @property {[]common.Entry} Entries - An array of Entry objects.
// @property {string} Cursor - A cursor for use with `ListFolderContinue`.
// @property {bool} HasMore - If true, then there are more entries available. Pass the cursor to
// ListFoldersContinue() to retrieve the rest.
type ListFoldersResponse struct {
	Entries []common.Entry `json:"entries"`
	Cursor  string         `json:"cursor"`
	HasMore bool           `json:"has_more"`
}

// It returns a list of directories in the root of the Dropbox account
func ListFoldersHandler(c *fiber.Ctx) error {

	auth, errO := utils.VerifyRoute(c, "dropbox")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()

	body, errM := json.Marshal(map[string]interface{}{
		"path":                                "",
		"recursive":                           false,
		"include_media_info":                  false,
		"include_deleted":                     false,
		"include_has_explicit_shared_members": false,
		"include_mounted_folders":             true,
		"include_non_downloadable_files":      false,
	})
	if errM != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	encode, _, err := service.Endpoints["ListFoldersEndpoint"].CallEncode([]interface{}{
		auth,
		string(body),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	var response ListFoldersResponse
	err = json.Unmarshal(encode, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	// Filter only directories
	var directories []common.Entry
	directories = append(directories, common.Entry{
		Tag:         ".folder",
		Name:        ".",
		PathLower:   "/",
		PathDisplay: "/",
	})
	for _, entry := range response.Entries {
		if entry.Tag == "folder" {
			directories = append(directories, entry)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   directories,
			"fields": []string{"name", "path_lower"},
		},
	})
}

// `ListFilesResponse` is a struct with three fields: `Entries`, `Cursor`, and `HasMore`.
//
// The `Entries` field is an array of `common.Entry` structs. The `Cursor` field is a string. The
// `HasMore` field is a boolean.
// @property {[]common.Entry} Entries - An array of entries.
// @property {string} Cursor - A string that can be used to retrieve the next page of results.
// @property {bool} HasMore - If true, then there are more entries available after this
type ListFilesResponse struct {
	Entries []common.Entry `json:"entries"`
	Cursor  string         `json:"cursor"`
	HasMore bool           `json:"has_more"`
}

// It takes a path from the query string, calls the Dropbox API to list all files and folders in that
// path, and returns only the files
func ListFilesHandler(c *fiber.Ctx) error {

	auth, errO := utils.VerifyRoute(c, "dropbox")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()
	path := c.Query("path")
	if path == "default" {
		path = ""
	}

	body, errM := json.Marshal(map[string]interface{}{
		"path":                                path,
		"recursive":                           false,
		"include_media_info":                  false,
		"include_deleted":                     false,
		"include_has_explicit_shared_members": false,
		"include_mounted_folders":             true,
		"include_non_downloadable_files":      true,
	})
	if errM != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	encode, _, err := service.Endpoints["ListFoldersEndpoint"].CallEncode([]interface{}{
		auth,
		string(body),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	var response ListFoldersResponse
	err = json.Unmarshal(encode, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	// List only files
	var files []common.Entry
	for _, entry := range response.Entries {
		if entry.Tag == "file" {
			files = append(files, entry)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   files,
			"fields": []string{"name", "id"},
		},
	})
}
