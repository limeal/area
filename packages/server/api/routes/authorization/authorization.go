package authorization

import (
	"area-server/authenticators"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/services"
	"area-server/store"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// `CreateAuthorizationForServiceRequest` is a struct with three fields: `AuthenticatorName`, `Code`,
// and `RedirectURI`.
//
// The `json:"authenticator"` tag tells the JSON encoder/decoder to use the field name `authenticator`
// when marshaling/unmarshaling.
//
// The `validate:"required"` tag tells the `validator` package to validate that the field is not empty.
// @property {string} AuthenticatorName - The name of the authenticator that you want to use to
// authenticate the user.
// @property {string} Code - The code that was returned from the authorization request.
// @property {string} RedirectURI - The redirect URI that was used to obtain the authorization code.
type CreateAuthorizationForServiceRequest struct {
	AuthenticatorName string `json:"authenticator" validate:"required"`
	Code              string `json:"code" validate:"required"`
	RedirectURI       string `json:"redirect_uri" validate:"required"`
}

// `AuthorizationMeta` is a struct with two fields, `Authenticator` and `Applets`.
//
// The `Authenticator` field is a string, and the `Applets` field is a slice of strings.
// @property {string} Authenticator - The name of the authenticator that will be used to authenticate
// the user.
// @property {[]string} Applets - A list of applets that are allowed to be used with this
// authenticator.
type AuthorizationMeta struct {
	Authenticator string   `json:"authenticator"`
	Applets       []string `json:"applets"`
}

// METHOD: GET
func GetAuthorizations(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)

	// Retrieve authorizations for the account
	var authorizations []models.Authorization
	if result := postgres.DB.Where(&models.Authorization{AccountUUID: account.UUID}).Find(&authorizations); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	var meta []AuthorizationMeta
	for _, authorization := range authorizations {
		var areas []models.Area
		if result := postgres.DB.Where(&models.Area{
			AuthorizationUUID: authorization.UUID,
		}).Find(&areas); result.Error != gorm.ErrRecordNotFound && result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}

		applets := make(map[uuid.UUID]string)
		for _, area := range areas {
			var applet models.Applet
			if result := postgres.DB.Where(&models.Applet{
				UUID: area.AppletUUID,
			}).First(&applet); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":  fiber.StatusInternalServerError,
					"error": "Internal server error",
				})
			}
			applets[applet.UUID] = applet.Name
		}

		var appletsUsed []string
		for _, value := range applets {
			appletsUsed = append(appletsUsed, value)
		}
		meta = append(meta, AuthorizationMeta{
			Authenticator: authorization.AuthService,
			Applets:       appletsUsed,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"authorizations": authorizations,
			"meta":           meta,
		},
	})
}

// METHOD: POST
// Body: service, code, redirect_uri

// It creates an authorization for a service
func CreateAuthorization(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)

	validate := validator.New()
	body := new(CreateAuthorizationForServiceRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Wrong body)",
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Bad Request (Invalid body)",
		})
	}

	// Get Service
	authenticator := authenticators.GetAuthenticator(body.AuthenticatorName)
	if authenticator == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Authenticator not found",
		})
	}

	fSR, fSE := authenticator.Authenticate([]interface{}{
		false,
		body.Code,
		body.RedirectURI,
	})

	if fSE != nil {
		return c.Status(fSE.StatusCode).JSON(fiber.Map{
			"code":  fSE.StatusCode,
			"error": fSE.ErrorDesc,
		})
	}

	if fSR == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Cannot create authorization for a service that does not use it",
		})
	}

	var tmpAuth models.Authorization
	if result := postgres.DB.Where(&models.Authorization{
		AuthService: authenticator.Name,
		AccessToken: fSR.Data["access_token"].(string),
	}).First(&tmpAuth); result.Error == nil && result.RowsAffected != 0 {

		if tmpAuth.AccountUUID == account.UUID {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"code": fiber.StatusOK,
				"data": fiber.Map{
					"message": "Already authorized !",
				},
			})
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "[Error] Cannot authorize a service for an account that is already linked to another account",
		})
	}

	// Retrieve authorizations for the account
	var authorizations []models.Authorization
	if result := postgres.DB.Where(&models.Authorization{
		AccountUUID: account.UUID,
		AuthService: authenticator.Name,
	}).Find(&authorizations); result.Error != nil || len(authorizations) != 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Service already authorized",
		})
	}

	var other []byte
	if fSR.Data["other"] != nil {
		var errO error
		other, errO = json.Marshal(fSR.Data["other"])
		if errO != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	}

	// Create authorization
	if result := postgres.DB.Create(&models.Authorization{
		UUID:         uuid.New(),
		AccountUUID:  account.UUID,
		Type:         fSR.Type,
		AuthService:  authenticator.Name,
		AccessToken:  fSR.Data["access_token"].(string),
		RefreshToken: fSR.Data["refresh_token"].(string),
		ExpireAt:     fSR.Data["expired_at"].(time.Time),
		Other:        other,
		Permanent:    false,
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": fiber.Map{
			"message": "Authorization created !",
		},
	})
}

// METHOD: DELETE
// Params: authService
func DeleteAuthorization(c *fiber.Ctx) error {

	authenticator := authenticators.GetAuthenticator(c.Params("name", ""))

	if authenticator == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "Authenticator not found",
		})
	}
	// Retrieve account from context
	account := c.Locals("account").(models.Account)

	// Retrieve authorizations for the account
	var authorizations []models.Authorization
	if result := postgres.DB.Where(&models.Authorization{
		AccountUUID: account.UUID,
		AuthService: authenticator.Name,
	}).Find(&authorizations); result.Error != nil {

		if result.Error != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error :>" + result.Error.Error(),
			})
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code": fiber.StatusNotFound,
			"data": fiber.Map{
				"message": "Authorization not found",
			},
		})
	}

	if len(authorizations) != 1 {
		return c.Status(fiber.StatusTeapot).JSON(fiber.Map{
			"code":  fiber.StatusTeapot,
			"error": "Internal server error :> " + strconv.Itoa(len(authorizations)) + " authorizations found",
		})
	}

	if authorizations[0].Permanent {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Cannot delete a permanent authorization",
		})
	}

	// Check if an area is using this authorization
	var areas []models.Area
	if result := postgres.DB.Where(&models.Area{
		AuthorizationUUID: authorizations[0].UUID,
	}).Find(&areas); result.Error != nil && result.Error != gorm.ErrRecordNotFound {

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error :>" + result.Error.Error(),
			})
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code": fiber.StatusForbidden,
			"data": fiber.Map{
				"message": "Service is used by an applet :> " + result.Error.Error(),
			},
		})
	}

	// Revoke authorization on the service
	if authenticator.AuthEndpoints.RevokeToken != nil {
		_, _, err := authenticator.AuthEndpoints.RevokeToken.CallEncode([]interface{}{authorizations[0]})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Cannot revoke token :>" + err.Error(),
			})
		}
	}

	// Stop trigger with applet uuid and delete applet
	for _, area := range areas {
		if err := store.StopTrigger(area.AppletUUID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Cannot stop trigger :>" + err.Error(),
			})
		}
		// Try to delete applet if it exists
		if result := postgres.DB.Where(&models.Applet{
			UUID: area.AppletUUID,
		}).Delete(&models.Applet{}); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Cannot delete applet :>" + result.Error.Error(),
			})
		}
	}

	// Delete authorization
	if result := postgres.DB.Delete(&authorizations[0]); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Cannot delete authorization :>" + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Authorization deleted !",
		},
	})
}

// It returns a list of services that the user is authorized to use
func GetAuthorizedServices(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)

	// Retrieve authorizations for the account
	var authorizations []models.Authorization
	if result := postgres.DB.Where(&models.Authorization{AccountUUID: account.UUID}).Find(&authorizations); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	result := map[string]bool{}
	for _, service := range services.List {
		name := service.Name
		result[name] = false
		if service.Authenticator == nil {
			result[name] = true
			continue
		}
		for _, authorization := range authorizations {
			if service.Authenticator.Name == authorization.AuthService {
				result[name] = true
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"services": result,
		},
	})
}
