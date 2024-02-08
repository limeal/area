package middlewares

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// It takes a token, validates it, and returns the claims if it's valid
func ValidateToken(token string) (map[string]interface{}, error) {
	if len(token) < 7 || token[:6] != "Bearer" {
		return nil, errors.New("Invalid token format")
	}

	claims, err := utils.ValidateJWT(token[7:])
	if err != nil || claims["accountID"] == nil {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

// It gets the token from the header, validates it, and gets the user from the token
func TokenMiddleware(c *fiber.Ctx) error {

	// 1. Get token from header
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "No token",
		})
	}

	// 2. Validate token
	claims, err := ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": err.Error(),
		})
	}

	// 3. Get user from token

	var user models.Account
	result := postgres.DB.Where("uuid = ?", claims["accountID"]).First(&user)

	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":  401,
			"error": "Token is not linked to an user",
		})
	}

	c.Locals("account", user)
	return c.Next()
}
