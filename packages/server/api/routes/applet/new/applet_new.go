package new

import (
	"area-server/classes/static"
	"area-server/classes/triggers"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/services"
	"area-server/store"
	"area-server/store/webhooks"
	"encoding/json"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// METHOD: GET
// Description: Check if user has a new applet and return component
func GetNewApplet(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)
	field := c.Params("field", "")

	var applet models.Applet
	// State (3 states): partial, complete, error
	if result := postgres.DB.Where(&models.Applet{AccountUUID: account.UUID, State: "partial"}).First(&applet); result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {

		// Create new applet
		applet = models.Applet{
			UUID:        uuid.New(),
			AccountUUID: account.UUID,
			State:       "partial",
			Status:      "stopped",
		}

		if result := postgres.DB.Create(&applet); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	}

	data := make(fiber.Map)

	var act models.Area
	if (field == "" || field == "action") && applet.Action != "" && postgres.DB.Where(&models.Area{AppletUUID: applet.UUID, Type: "action"}).First(&act).RowsAffected != 0 {
		data["action"] = act
	}

	if field == "" || field == "reactions" {
		var reactions []models.Area
		if result := postgres.DB.Where(&models.Area{AppletUUID: applet.UUID, Type: "reaction"}).Find(&reactions); result.RowsAffected != 0 {
			data["reactions"] = reactions
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": data,
	})
}

// `AddStateToNewAppletRequest` is a struct with fields `Service`, `AreaType`, `AreaItem`, and
// `AreaItemSettings`.
//
// The `validate` tag is used to specify validation rules for the field. In this case, the `Service`
// field is required, and the `AreaType` field must be either `action` or `reaction`.
//
// The `json` tag is used to specify the name of the field in the JSON request.
// @property {string} Service - The name of the service you want to add to the applet.
// @property {string} AreaType - This is the type of applet you're creating. It can be either an action
// or a reaction.
// @property {string} AreaItem - The name of the item you want to add to the applet.
// @property AreaItemSettings - This is a map of settings that are specific to the area item. For
// example, if the area item is "send_email", then the area settings would be the email address to send
// the email to.
type AddStateToNewAppletRequest struct {
	// Type can be "action" or "reaction"
	Service          string                 `json:"service" validate:"required"`
	AreaType         string                 `json:"area_type" validate:"required,oneof=action reaction"`
	AreaItem         string                 `json:"area_item" validate:"required"`
	AreaItemSettings map[string]interface{} `json:"area_settings"`
}

// METHOD: PUT
// Description: Add state to new applet
func AddStateToNewApplet(c *fiber.Ctx) error {

	validate := validator.New()
	applet := c.Locals("applet").(models.Applet)
	account := c.Locals("account").(models.Account)

	// Read body
	body := new(AddStateToNewAppletRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid body",
		})
	}

	// Validate body
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid body",
		})
	}

	store, err := json.Marshal(body.AreaItemSettings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	// Check if a area of the same type already exists
	if body.AreaType == "action" && applet.Action != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Applet action already linked",
		})
	}

	service := services.GetServiceByName(body.Service)
	if service == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "Service not found",
		})
	}

	// Get Authorization
	var authorization models.Authorization
	authorization.UUID = uuid.Nil
	if service.Authenticator != nil {
		if result := postgres.DB.Where(&models.Authorization{
			AccountUUID: account.UUID,
			AuthService: service.Authenticator.Name,
		}).First(&authorization); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":  fiber.StatusNotFound,
				"error": "Authorization not found :>" + result.Error.Error(),
			})
		}
	}

	// Check if action/reaction exists
	var areaItem *static.ServiceArea
	if body.AreaType == "action" {
		areaItem = service.GetActionByName(body.AreaItem)
	} else {
		areaItem = service.GetReactionByName(body.AreaItem)
	}

	if areaItem == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":  fiber.StatusNotFound,
			"error": "Area item not found",
		})
	}

	if service.Authenticator != nil {
		_, err = service.Authenticator.RefreshToken(&authorization)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":  fiber.StatusUnauthorized,
				"error": "Authorization expired",
			})
		}
	}

	// Foreach settings, if value is not empty, check if it's valid
	var invalidKeys map[string]bool

	if authorization.UUID != uuid.Nil {
		invalidKeys = areaItem.Validate(&authorization, service, body.AreaItemSettings)
	} else {
		invalidKeys = areaItem.Validate(nil, service, body.AreaItemSettings)
	}

	if invalidKeys != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"code": fiber.StatusNotAcceptable,
			"data": invalidKeys,
		})
	}

	area := &models.Area{
		UUID:              uuid.New(),
		AppletUUID:        applet.UUID,
		AuthorizationUUID: authorization.UUID,
		Type:              body.AreaType,
		Service:           body.Service,
		Name:              body.AreaItem,
		Store:             store,
	}

	// Create AREA
	if result := postgres.DB.Create(&area); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	if body.AreaType == "action" {
		applet.Action = area.Service + ";" + area.Name
	}

	if result := postgres.DB.Save(&applet); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet updated !",
		},
	})
}

// `SubmitNewAppletRequest` is a struct with three fields, `Name`, `Description`, and `Public`.
//
// The `Name` field is a string, and it's required.
//
// The `Description` field is a string, and it's required.
//
// The `Public` field is a bool, and it's not required.
// @property {string} Name - The name of the applet.
// @property {string} Description - The description of the applet.
// @property {bool} Public - Whether the applet is public or not.
type SubmitNewAppletRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Public      bool   `json:"public"`
}

// METHOD: POST
// Description: Submit new applet
func SubmitNewApplet(c *fiber.Ctx) error {

	// Read body
	validate := validator.New()
	body := new(SubmitNewAppletRequest)
	account := c.Locals("account").(models.Account)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid body",
		})
	}

	// Validate body
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": err.Error(),
		})
	}

	applet := c.Locals("applet").(models.Applet)

	var reactions []models.Area
	if result := postgres.DB.Where(&models.Area{
		AppletUUID: applet.UUID,
		Type:       "reaction",
	}).Find(&reactions); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	// Check if applet is complete (has 1 action and 1 reaction)
	if applet.Action == "" || len(reactions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Applet is not complete",
		})
	}

	// Check if applet with same name already exists
	if result := postgres.DB.Where(&models.Applet{
		AccountUUID: account.UUID,
		Name:        body.Name,
	}).First(&models.Applet{}); result.Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "You have already an applet with this name",
		})
	}

	// Check if applet action is webhook
	if applet.Action == "webhook;applet_triggered" {
		webhooks.AddWebhook(applet.UUID.String())
	}

	// Create trigger
	tr, err := triggers.CreateTrigger(&applet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}
	store.AddTrigger(tr, true)

	// Update applet
	if result := postgres.DB.Model(&applet).Updates(models.Applet{
		Name:        body.Name,
		Description: body.Description,
		Public:      body.Public || false,
		Status:      "running",
		State:       "complete",
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet submitted !",
		},
	})
}

// METHOD: DELETE
// Query: type=action|reaction
func DeleteStateNewApplet(c *fiber.Ctx) error {

	applet := c.Locals("applet").(models.Applet)
	atype := c.Query("type")
	number := c.Query("number", "")

	if atype != "action" && atype != "reaction" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid query",
		})
	}

	if atype == "action" {

		if applet.Action == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":  fiber.StatusNotFound,
				"error": "Action not found",
			})
		}

		if result := postgres.DB.Where(&models.Applet{
			AccountUUID: applet.AccountUUID,
			State:       "partial",
		}).Delete(&models.Applet{}); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	}

	if atype == "reaction" {
		var reactions []models.Area
		if result := postgres.DB.Where(&models.Area{
			AppletUUID: applet.UUID,
			Type:       "reaction",
		}).Order("created_at asc").
			Find(&reactions); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}

		if number != "" {
			num, err := strconv.Atoi(number)
			if err != nil || num >= len(reactions) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"code":  fiber.StatusBadRequest,
					"error": "Invalid query",
				})
			}

			if result := postgres.DB.Where(reactions[num].UUID).Delete(&models.Area{}); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":  fiber.StatusInternalServerError,
					"error": "Internal server error",
				})
			}
		} else {
			if result := postgres.DB.Where(&models.Area{
				AppletUUID: applet.UUID,
				Type:       "reaction",
			}).Delete(&models.Area{}); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":  fiber.StatusInternalServerError,
					"error": "Internal server error",
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": atype + " deleted !",
		},
	})
}
