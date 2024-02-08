package user

import (
	"area-server/db/postgres/models"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// It takes a multipart form, saves the file to the `./avatars` directory, and returns a success
// message
func UpdateAvatar(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)

	mult, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid multipart form",
		})
	}

	files := mult.File["avatar"]

	for _, file := range files {
		if err := c.SaveFile(file, fmt.Sprintf("./avatars/%s", account.UUID)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Error while saving file",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Avatar updated",
		},
	})
}

// It checks if the avatar exists, if it does, it returns the avatar's URI, if it doesn't, it returns
// the default avatar's URI
func GetAvatar(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)
	base := c.Query("base", "http://localhost:8080")

	fileName := fmt.Sprintf("./avatars/%s", account.UUID)
	_, error := os.Stat(fileName)
	// Check if avatar exists
	// If not, return default avatar
	if os.IsNotExist(error) {
		return c.Status(200).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"data": fiber.Map{
				"uri": "https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png",
			},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"uri": fmt.Sprintf(base+"/avatars/%s", account.UUID),
		},
	})
}
