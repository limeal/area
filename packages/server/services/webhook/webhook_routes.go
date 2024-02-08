package webhook

import (
	"area-server/classes/static"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/store/webhooks"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of static.ServiceRoute, which is a struct that contains the endpoint, the
// handler, the method, and whether or not the route needs authentication
func WebhookRoute() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/:name",
			Handler:  webhookHandler,
			Method:   "POST",
			NeedAuth: true,
		},
		{
			Endpoint: "/:email/:name",
			Handler:  webhookHandlerWithoutAuth,
			Method:   "POST",
			NeedAuth: false,
		},
		{
			Endpoint: "/history",
			Handler:  historyHandler,
			Method:   "GET",
			NeedAuth: false,
		},
	}
}

// --------------------- Routes ---------------------

// It receives a webhook request, checks if the applet exists, and writes the request body to a channel
func webhookHandler(c *fiber.Ctx) error {
	appletName := c.Params("name", "")
	account := c.Locals("account").(models.Account)

	// Check if the applet exists
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{
		AccountUUID: account.UUID,
		Name:        appletName,
	}).First(&applet); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "Applet not found",
		})
	}

	fmt.Println("Webhook received for applet: " + appletName)
	// Write to channel
	if err := webhooks.WriteToWebhook(account.Email, applet.UUID.String(), c.Body()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Webhook received",
		},
	})
}

// It checks if the user and applet exists, and if so, it writes the webhook to a channel
func webhookHandlerWithoutAuth(c *fiber.Ctx) error {
	userEmail := c.Params("email", "")
	appletName := c.Params("name", "")

	// Check if user exist
	var account models.Account
	if result := postgres.DB.Where(&models.Account{
		Email: userEmail,
	}).First(&account); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "User not found",
		})
	}

	// Check if the applet exists (only public is allowed in this case)
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{
		AccountUUID: account.UUID,
		Name:        appletName,
		Public:      true,
	}).First(&applet); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "Applet not found",
		})
	}

	// Write to channel
	if err := webhooks.WriteToWebhook(userEmail, applet.UUID.String(), c.Body()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Webhook received",
		},
	})
}

// It prints the history of the webhook
func historyHandler(c *fiber.Ctx) error {
	// Print the history of the webhook

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"history": webhooks.History,
		},
	})
}
