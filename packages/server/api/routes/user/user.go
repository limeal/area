package user

import (
	"area-server/authenticators"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	sessionr "area-server/store"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// NEED AUTHENTICATION

// > It gets the account from the context and returns it as JSON
func GetUser(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"account": account,
		},
	})
}

// It destroys the session of the user and returns a success message
func LogoutUser(c *fiber.Ctx) error {

	if c.Locals("session") == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "No session found",
		})
	}

	session := c.Locals("session").(*session.Session)
	if !sessionr.DestroySession(session) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "User logged out",
		},
	})
}

// `UpdateUserBody` is a struct with a single field, `Username`, which is a string.
//
// The `json:"username"` part is called a struct tag. It tells the JSON encoder/decoder how to map the
// field to JSON. In this case, it tells the encoder to use the field name `username` when encoding to
// JSON.
//
// The `json:"username"` part is called a struct tag. It tells the JSON encoder/decoder how to map the
// field to JSON. In this case, it tells the encoder to
// @property {string} Username - The username of the user to update.
type UpdateUserBody struct {
	Username string `json:"username"`
}

// We validate the body, retrieve the account from the context, and update the account
func UpdateUser(c *fiber.Ctx) error {

	validate := validator.New()
	body := new(UpdateUserBody)

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Wrong Body)",
		})
	}

	// Validate body
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Invalid Body)",
		})
	}

	// Retrieve account from context
	account := c.Locals("account").(models.Account)

	// Update the account
	if result := postgres.DB.Model(&account).Updates(models.Account{
		Username: body.Username,
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "User updated",
		},
	})
}

// It deletes the user's account, revokes all of their tokens, and destroys their session
func DeleteUser(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)

	var Authorization []models.Authorization
	if result := postgres.DB.Where(&models.Authorization{
		AccountUUID: account.UUID,
	}).Find(&Authorization); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	for _, authorization := range Authorization {
		authenticator := authenticators.GetAuthenticator(authorization.AuthService)
		if authenticator == nil || authenticator.AuthEndpoints.RevokeToken == nil {
			continue
		}
		_, _, err := authenticator.AuthEndpoints.RevokeToken.CallEncode([]interface{}{authorization})
		if err != nil {
			fmt.Println("Could not revoke token for: ", authorization.AuthService)
		}
	}

	if result := postgres.DB.Delete(&account); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	if c.Locals("session") != nil {
		if !sessionr.DestroySession(c.Locals("session").(*session.Session)) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "User deleted",
		},
	})
}

// Redirections

// It redirects the user to the applet route with the query string `public=false`
func GetUserApplets(c *fiber.Ctx) error {
	return c.RedirectToRoute("applet", fiber.Map{
		"queries": map[string]string{
			"public": "false",
		},
	})
}
