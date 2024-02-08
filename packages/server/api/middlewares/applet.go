package middlewares

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"

	"github.com/gofiber/fiber/v2"
)

// > This function will get the applet from the database and add it to the context
func NewAppletMiddleware(c *fiber.Ctx) error {

	// Need to add session / token middleware before this one
	account := c.Locals("account").(models.Account)

	// 3. Get applet from database
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{AccountUUID: account.UUID, State: "partial"}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "No Applet found",
		})
	}

	c.Locals("applet", applet)
	return c.Next()
}
