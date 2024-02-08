package store

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"

	"github.com/gofiber/fiber/v2"
)

// Get all applets that are public and in a complete state
func GetStoreApplets(c *fiber.Ctx) error {

	// 1. Get all applets public only

	var applets []models.Applet
	if result := postgres.DB.Where(&models.Applet{State: "complete", Public: true}).Find(&applets); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applets",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"applets": applets,
		},
	})
}
