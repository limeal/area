package authenticators

import (
	"area-server/authenticators"

	"github.com/gofiber/fiber/v2"
)

// `GetAuthenticators` is a function that takes a `fiber.Ctx` and returns an `error`
func GetAuthenticators(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(authenticators.List)
}

// "Get the authenticator with the given name and return it as JSON."
//
// The first thing we do is get the name of the authenticator from the URL. We then loop through all
// the authenticators and check if the name matches. If it does, we return the authenticator as JSON.
// If it doesn't, we return a 404 error
func GetAuthenticator(c *fiber.Ctx) error {
	name := c.Params("name")

	for _, authenticator := range authenticators.List {
		if authenticator.Name == name {
			return c.Status(fiber.StatusOK).JSON(authenticator)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Authenticator not found",
	})
}
