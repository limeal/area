package services

import (
	sservices "area-server/services"

	"github.com/gofiber/fiber/v2"
)

// > The function `GetServices` takes a `fiber.Ctx` object as an argument and returns an error
func GetServices(c *fiber.Ctx) error {
	return c.Status(200).JSON(sservices.List)
}

// It takes the service name from the URL, loops through the list of services, and if it finds a match,
// returns the service. If it doesn't find a match, it returns a 404 error
func GetService(c *fiber.Ctx) error {
	service := c.Params("service")

	for _, s := range sservices.List {
		if s.Name == service {
			return c.Status(200).JSON(s)
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Service not found",
	})
}

// It takes the service name from the URL, loops through the list of services, and if it finds a match,
// it returns the list of actions for that service
func GetServiceActions(c *fiber.Ctx) error {
	service := c.Params("service")

	for _, s := range sservices.List {
		if s.Name == service {
			return c.Status(200).JSON(s.Actions)
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Service not found",
	})
}

// It gets the service and action from the URL, loops through the services and actions, and returns the
// action if it exists
func GetServiceAction(c *fiber.Ctx) error {
	service := c.Params("service")
	action := c.Params("action")

	for _, s := range sservices.List {
		if s.Name == service {
			for _, a := range s.Actions {
				if a.Name == action {
					return c.Status(200).JSON(a)
				}
			}
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Action not found",
	})
}

// It gets the service name from the URL, loops through the list of services, and if it finds a match,
// it returns the service's reactions
func GetServiceReactions(c *fiber.Ctx) error {
	service := c.Params("service")

	for _, s := range sservices.List {
		if s.Name == service {
			return c.Status(200).JSON(s.Reactions)
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Service not found",
	})
}

// It gets the service and reaction from the URL, loops through the services and reactions, and returns
// the reaction if it exists
func GetServiceReaction(c *fiber.Ctx) error {
	service := c.Params("service")
	reaction := c.Params("reaction")

	for _, s := range sservices.List {
		if s.Name == service {
			for _, r := range s.Reactions {
				if r.Name == reaction {
					return c.Status(200).JSON(r)
				}
			}
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Reaction not found",
	})
}

// It takes the service name from the URL, then loops through the list of services and returns the
// routes for the service if it exists
func GetApiEndpoints(c *fiber.Ctx) error {
	service := c.Params("service")

	for _, s := range sservices.List {
		if s.Name == service {
			return c.Status(200).JSON(s.Routes)
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Service not found",
	})
}
