package spotify

import (
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of `static.ServiceRoute`s
func SpotifyRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/playlists",
			Handler:  GetPlaylistsRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/devices",
			Handler:  GetUserDevices,
			Method:   "GET",
			NeedAuth: true,
		},
	}
}

// ---------------------- Spotify Routes ----------------------

// `GetUserPlaylistsResponse` is a struct with a field `Items` which is a slice of structs with fields
// `ID`, `Name`, `Collaborative`, and `Owner` which is a struct with a field `ID`.
// @property {[]struct {
// 		ID            string `json:"id"`
// 		Name          string `json:"name"`
// 		Collaborative bool   `json:"collaborative"`
// 		Owner         struct {
// 			ID string `json:"id"`
// 		} `json:"owner"`
// 	}} Items - An array of playlist objects.
type GetUserPlaylistsResponse struct {
	Items []struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Collaborative bool   `json:"collaborative"`
		Owner         struct {
			ID string `json:"id"`
		} `json:"owner"`
	}
}

// It gets the user's playlists and returns them
func GetPlaylistsRoute(c *fiber.Ctx) error {

	auth, errO := utils.VerifyRoute(c, "spotify")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}
	// Check if playlist is collaborative or if user is owner (queryParams modify)
	service := Descriptor()

	// Get List of playlists
	encode, _, err := service.Endpoints["GetUserPlaylistsEndpoint"].CallEncode([]interface{}{
		auth,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	var playlists GetUserPlaylistsResponse
	err = json.Unmarshal(encode, &playlists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   playlists.Items,
			"fields": []string{"name", "id"},
		},
	})
}

// `GetUserDevicesResponse` is a struct with a field `Devices` which is an array of structs with fields
// `ID` and `Name`.
// @property {[]struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 	}} Devices - An array of devices that the user has.
type GetUserDevicesResponse struct {
	Devices []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"devices"`
}

// It gets the list of devices that the user has connected to their Spotify account
func GetUserDevices(c *fiber.Ctx) error {

	auth, errO := utils.VerifyRoute(c, "spotify")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()

	// Get List of devices
	encode, _, err := service.Endpoints["GetUserAvailableDevicesEndpoint"].CallEncode([]interface{}{
		auth,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	deviceResp := GetUserDevicesResponse{}
	err = json.Unmarshal(encode, &deviceResp.Devices)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   deviceResp.Devices,
			"fields": []string{"name", "id"},
		},
	})
}
