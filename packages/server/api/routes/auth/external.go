package auth

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/config"
	pg "area-server/db/postgres"
	models "area-server/db/postgres/models"
	session "area-server/store"
	"area-server/utils"
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// `ExternalAuthRequest` is a struct with three fields, `Authenticator`, `Code`, and `RedirectURI`, all
// of which are strings.
// @property {string} Authenticator - The name of the authenticator you want to use.
// @property {string} Code - The code that was returned from the external authentication provider.
// @property {string} RedirectURI - The redirect URI that was used to obtain the code.
type ExternalAuthRequest struct {
	Authenticator string `validate:"required" json:"authenticator"`
	Code          string `validate:"required" json:"code"`
	RedirectURI   string `validate:"required" json:"redirect_uri"`
}

// It creates a new account and authorization for the user
func newAccount(fSR *static.AuthenticateResponse, externalAuthReq *ExternalAuthRequest) (*models.Account, error) {
	account := &models.Account{
		UUID:          uuid.New(),
		Email:         fSR.Data["email"].(string),
		Authenticator: externalAuthReq.Authenticator,
		Username:      utils.GenerateNameForUser(),
	}

	if result := pg.DB.Create(account); result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	var other []byte
	if fSR.Data["other"] != nil {
		var errO error
		other, errO = json.Marshal(fSR.Data["other"])
		if errO != nil {
			return nil, errO
		}
	}

	// Create authorization
	if err := pg.DB.Create(&models.Authorization{
		UUID:         uuid.New(),
		AccountUUID:  account.UUID,
		Type:         "oauth2",
		AuthService:  externalAuthReq.Authenticator,
		AccessToken:  fSR.Data["access_token"].(string),
		RefreshToken: fSR.Data["refresh_token"].(string),
		ExpireAt:     fSR.Data["expired_at"].(time.Time),
		Other:        other,
		Permanent:    true,
	}).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// External Body example:
/*
 * GitHub Login
 * Service: "github"
 * Code: "code"
 * RedirectURI: "http://localhost:3000/auth/github"
 *
 */

// It takes a code, authenticator, and redirect URI, and returns a token or session
func ExternalAuth(c *fiber.Ctx) error {

	// Token Mode: Remove this
	sess, serr := session.SessionStore.Get(c)
	if serr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error (Store)",
		})
	}

	validate := validator.New()
	externalAuthReq := new(ExternalAuthRequest)

	if err := c.BodyParser(externalAuthReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "[Error] Bad Request (JSON)",
		})
	}

	if err := validate.Struct(externalAuthReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "[Error] Bad Request (JSON)",
		})
	}

	// Get Authenticator
	authenticator := authenticators.GetAuthenticator(externalAuthReq.Authenticator)
	if authenticator == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "[Error] Authenticator not found",
		})
	}

	if authenticator == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Service does not use authorization",
		})
	}

	// Validate service + Retrieve Email
	fSR, fSE := authenticator.Authenticate([]interface{}{
		true,
		externalAuthReq.Code,
		externalAuthReq.RedirectURI,
	})

	if fSE != nil {
		return c.Status(fSE.StatusCode).JSON(fiber.Map{
			"code":  fSE.StatusCode,
			"error": fSE.ErrorDesc,
		})
	}

	// 4. Check if account exists, if not, return 401
	var Account models.Account
	var status int
	var message string
	var accountUUID uuid.UUID

	if result := pg.DB.Where(&models.Account{Email: fSR.Data["email"].(string)}).First(&Account); result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		// Create account

		// Check if access token already exists
		if result := pg.DB.Where(&models.Authorization{
			AuthService: externalAuthReq.Authenticator,
			AccessToken: fSR.Data["access_token"].(string),
		}).First(&models.Authorization{}); result.Error == nil && result.RowsAffected != 0 {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"code":  fiber.StatusNotAcceptable,
				"error": "[Error] Duplicate access token",
			})
		}

		account, err := newAccount(fSR, externalAuthReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
		status, message, accountUUID = 201, "Account Created", account.UUID
	} else {
		if Account.Authenticator != externalAuthReq.Authenticator {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"code":  fiber.StatusNotAcceptable,
				"error": "Account already exist with this email !",
			})
		}

		status, message, accountUUID = 200, "Account Logged", Account.UUID
	}

	data := make(map[string]interface{})

	data["message"] = message
	if config.CFG.Mode == config.Session {
		session.DestroySession(sess)
		err := session.GenerateSession(sess, accountUUID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
	} else {
		token, err := utils.GenerateJWT(fiber.Map{
			"accountID": accountUUID,
		}, config.CFG.TokenDuration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
		data["token"] = token
	}

	return c.Status(status).JSON(fiber.Map{
		"code": status,
		"data": data,
	})

}
