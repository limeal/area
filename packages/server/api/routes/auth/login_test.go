package auth

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// LoginRequestTest is a struct that contains two fields, Email and EncodedPassword, both of which are
// strings.
// @property {string} Email - The email address of the user.
// @property {string} EncodedPassword - The encoded password that the user entered.
type LoginRequestTest struct {
	Email           string `validate:"required,email,min=6,max=32" json:"email"`
	EncodedPassword string `validate:"required" json:"encoded_password"`
}

// We're creating a new Fiber app, and then we're adding a new route to it
func TestLogin(t *testing.T) {

	app := fiber.New()
	defer app.Shutdown()

	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		validate := validator.New()
		loginReq := new(LoginRequestTest)

		// Parse JSON
		if err := c.BodyParser(loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":  fiber.StatusBadRequest,
				"error": "Bad Request (Wrong Body)",
			})
		}

		// Validate
		if err := validate.Struct(loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":  fiber.StatusBadRequest,
				"error": "Bad Request (Invalid Body)",
			})
		}

		return c.Status(fiber.StatusOK).SendString("OK")
	})

	tests := []fiber.Map{
		{
			"email":            "loul@example.com",
			"encoded_password": "12344521",
			"status":           200,
		},
		{
			"email":            "paul",
			"encoded_password": "12344521",
			"status":           400,
		},
		{
			"email":  "test",
			"status": 400,
		},
		{
			"password": "test",
			"status":   400,
		},
		{
			"status": 400,
		},
	}

	for _, tt := range tests {
		body := map[string]interface{}{
			"email":            tt["email"],
			"encoded_password": tt["encoded_password"],
		}

		str, err := json.Marshal(body)
		if err != nil {
			continue
		}

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(str))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)

		assert.Equalf(t, tt["status"], resp.StatusCode, "Status code should be %d", tt["status"])
	}
}
