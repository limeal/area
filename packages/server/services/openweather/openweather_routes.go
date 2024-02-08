package openweather

import (
	"area-server/classes/static"
	"area-server/services/openweather/common"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of `static.ServiceRoute` structs
func OpenWeatherRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/stations",
			Method:   "GET",
			Handler:  GetStationsHandler,
			NeedAuth: true,
		},
	}
}

// --------------------- Routes ---------------------

// It calls the `GetStationsEndpoint` endpoint, unmarshals the response into a `[]common.Station` and
// returns the result as a JSON response
func GetStationsHandler(c *fiber.Ctx) error {

	service := Descriptor()
	var stations []common.Station
	encode, _, err := service.Endpoints["GetStationsEndpoint"].CallEncode([]interface{}{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	if err := json.Unmarshal(encode, &stations); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   stations,
			"fields": []string{"name", "id"},
		},
	})
}
