package youtube

import (
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of `static.ServiceRoute`s
func YoutubeRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/playlists",
			Method:   "GET",
			Handler:  GetAllPlaylists,
			NeedAuth: true,
		},
	}
}

// --------------------- Routes ---------------------

// It gets all the playlists of the user
func GetAllPlaylists(c *fiber.Ctx) error {

	auth0, err := utils.VerifyRoute(c, "google")
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()

	// Get all playlists
	encode, _, err := service.Endpoints["GetAllPlaylistsEndpoint"].CallEncode([]interface{}{
		auth0,
		map[string]string{
			"part":       "snippet",
			"mine":       "true",
			"maxResults": "20",
		},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	playlists := common.YoutubePlaylistsResponse{}
	if err := json.Unmarshal(encode, &playlists); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   playlists.Items,
			"fields": []string{"snippet:title", "id"},
		},
	})
}
