package auth

import (
	"area-server/config"
	pg "area-server/db/postgres"
	models "area-server/db/postgres/models"
	session "area-server/store"
	"area-server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LoginRequest is a struct that contains two fields, Email and EncodedPassword, both of which are
// strings.
// @property {string} Email - The email address of the user.
// @property {string} EncodedPassword - The encoded password is the password that the user enters in
// the login form.
type LoginRequest struct {
	Email           string `validate:"required,email,min=6,max=32" json:"email"`
	EncodedPassword string `validate:"required,min=10" json:"encoded_password"`
}

/*
 *
 * Description: Route to login
 * Method: POST
 * Login Call example:
 * Body:
 *		Email: <email>
 *		EncodedPassword: <password>
 *
 * @Return 200: Login successful
 */
func Login(c *fiber.Ctx) error {

	sess, serr := session.SessionStore.Get(c)
	if serr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error (Store)",
		})
	}
	validate := validator.New()
	loginReq := new(LoginRequest)

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

	// Check if account already exist
	var Account models.Account
	if result := pg.DB.Where(&models.Account{Email: loginReq.Email, Authenticator: "@local"}).First(&Account); result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"code":  fiber.StatusNotAcceptable,
			"error": "Account not found",
		})
	}

	// Check if password is correct
	if loginReq.EncodedPassword != Account.Password {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"code":  fiber.StatusNotAcceptable,
			"error": "Account not found",
		})
	}

	data := make(map[string]interface{})

	data["message"] = "Login successful !"
	if config.CFG.Mode == config.Session {
		session.DestroySession(sess)
		err := session.GenerateSession(sess, Account.UUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
	} else {
		token, err := utils.GenerateJWT(fiber.Map{
			"accountID": Account.UUID,
		}, config.CFG.TokenDuration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
		data["token"] = token
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": data,
	})
}
