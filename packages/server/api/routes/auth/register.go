package auth

import (
	"area-server/config"
	pg "area-server/db/postgres"
	models "area-server/db/postgres/models"
	session "area-server/store"
	utils "area-server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Validator

// `RegisterRequest` is a struct that has two fields, `Email` and `EncodedPassword`, both of which are
// strings.
//
// The `validate` tag is what we're interested in. It's a string that contains a comma-separated list
// of validation rules. The rules are applied to the field they're attached to.
//
// The `required` rule means that the field must be present. The `email` rule means that the field must
// be a valid email address. The `min` and `max` rules mean that the field must be at least `min
// @property {string} Email - The email address of the user.
// @property {string} EncodedPassword - The encoded password is the password that the user will use to
// login.
type RegisterRequest struct {
	Email           string `validate:"required,email,min=6,max=32" json:"email"`
	EncodedPassword string `validate:"required,min=10" json:"encoded_password"`
}

// Handler
func Register(c *fiber.Ctx) error {

	sess, serr := session.SessionStore.Get(c)
	if serr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error (Store)",
		})
	}
	// Pre-Validation
	validate := validator.New()
	registerReq := new(RegisterRequest)

	// Parse JSON
	if err := c.BodyParser(registerReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Wrong Body)",
		})
	}

	// Validate
	if err := validate.Struct(registerReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Invalid Body)",
		})
	}

	// Check if account already exist
	var Account models.Account
	if result := pg.DB.Where(&models.Account{Email: registerReq.Email}).First(&Account); result.Error != gorm.ErrRecordNotFound || result.RowsAffected != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Account already exist",
		})
	}

	// Create account
	account := models.Account{
		UUID:          uuid.New(),
		Authenticator: "@local",
		Email:         registerReq.Email,
		Username:      utils.GenerateNameForUser(),
		Password:      registerReq.EncodedPassword,
	}

	if result := pg.DB.Create(&account); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	data := make(map[string]interface{})

	data["message"] = "Account created !"
	if config.CFG.Mode == config.Session {
		session.DestroySession(sess)
		// Generate new session
		err := session.GenerateSession(sess, account.UUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

	} else {
		token, err := utils.GenerateJWT(fiber.Map{
			"accountID": account.UUID,
		}, config.CFG.TokenDuration)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		data["token"] = token
	}

	// Return token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})
}
