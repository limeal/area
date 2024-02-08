package middlewares

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	session "area-server/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// It retrieves the session from the session store, checks if the account is linked to the session, and
// if it is, it retrieves the account from the database and stores it in the context
func SessionMiddleware(c *fiber.Ctx) error {

	sess, serr := session.SessionStore.Get(c)
	if serr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid session",
		})
	}

	accountID := sess.Get("account")
	if accountID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Session insn't valid (no account linked)",
		})
	}

	accountUUID, err := uuid.Parse(accountID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	// Retrieve account from database
	var Account models.Account

	if result := postgres.DB.Where(&models.Account{UUID: accountUUID}).First(&Account); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Account not found",
		})
	}

	c.Locals("account", Account)
	c.Locals("session", sess)
	return c.Next()
}
