package applet

import (
	"area-server/store"

	"area-server/db/postgres"
	models "area-server/db/postgres/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NEED AUTHENTICATION

// METHOD: GET
func GetApplets(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)

	var applets []models.Applet
	if result := postgres.DB.Where(&models.Applet{AccountUUID: account.UUID, State: "complete"}).Find(&applets); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applets",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"applets": applets,
		},
	})
}

// With Params -> action_id

// METHOD: Get
func GetApplet(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{AccountUUID: account.UUID, UUID: appletId, State: "complete"}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applet found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"applet": applet,
		},
	})
}

// Get all reactions for a given applet
func GetAppletReactions(c *fiber.Ctx) error {

	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	var reactions []models.Area
	if result := postgres.DB.Where(&models.Area{AppletUUID: appletId, Type: "reaction"}).Order("created_at asc").Find(&reactions); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"data": fiber.Map{
				"reactions": []models.Area{},
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"reactions": reactions,
		},
	})
}

// `UpdateAppletRequest` is a struct with fields `Service`, `AreaType`, `AreaItem`, and
// `AreaItemSettings`.
//
// The `validate` tag is used to specify validation rules for the field. In this case, the `Service`
// field is required, and the `AreaType` field must be either `action` or `reaction`.
//
// The `json` tag is used to specify the name of the field in the JSON request.
// @property {string} Service - The name of the service you want to update.
// @property {string} AreaType - The type of area the applet is for. This can be either "action" or
// "reaction".
// @property {string} AreaItem - The name of the item you want to update.
// @property AreaItemSettings - This is a map of settings that are specific to the area item. For
// example, if the area item is a webhook, then the area settings would be the URL of the webhook.
type UpdateAppletRequest struct {
	Service          string                 `json:"service" validate:"required"`
	AreaType         string                 `json:"area_type" validate:"required,oneof=action reaction"`
	AreaItem         string                 `json:"area_item" validate:"required"`
	AreaItemSettings map[string]interface{} `json:"area_settings"`
}

// METHOD: PUT
/* func UpdateApplet(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	// Parse body
	validate := validator.New()
	body := new(UpdateAppletRequest)

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

	// Get area of type action/reaction linked to applet
	var area models.Area
	if result := postgres.DB.Where(&models.Area{AppletUUID: appletId, Type: body.AreaType}).First(&area); result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Cannot find area of type " + body.AreaType,
		})
	}

	if area.Service == body.Service && area.Name == body.AreaItem {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code":  fiber.StatusConflict,
			"error": "Area already exists",
		})
	}

	bstore, err := json.Marshal(body.AreaItemSettings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
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
				"error": "Authorization not found",
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

	// Update area
	if result := postgres.DB.Model(&area).Updates(&models.Area{
		AuthorizationUUID: authorization.UUID,
		Service:           body.Service,
		Name:              body.AreaItem,
		Store:             bstore,
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	// Update Context (Only if trigger is running)
	store.UpdateTrigger(appletId, body.AreaType)

	// Update applet
	if body.AreaType == "action" {
		if result := postgres.DB.Model(&models.Applet{UUID: appletId}).Updates(&models.Applet{
			ActionLinked: body.Service + ";" + body.AreaItem,
		}); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	} else {
		if result := postgres.DB.Model(&models.Applet{UUID: appletId}).Updates(&models.Applet{
			ReactionLinked: body.Service + ";" + body.AreaItem,
		}); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Area updated",
		},
	})
} */

// METHOD: PATCH
func UpdateAppletActivity(c *fiber.Ctx) error {
	activity := c.Query("active", "")
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	if activity != "true" && activity != "false" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid activity",
		})
	}

	// Get applet from database
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{
		UUID:   appletId,
		Status: "running",
	}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applet found / Applet not running",
		})
	}

	if (activity == "true" && applet.Active == true) || (activity == "false" && applet.Active == false) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code":  fiber.StatusConflict,
			"error": "Applet already in this state",
		})
	}

	// Set applet activity
	if activity == "true" {
		store.OpenTrigger(applet.UUID)
	} else {
		store.CloseTrigger(applet.UUID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet activity updated",
		},
	})
}

// METHOD: DELETE
func DeleteApplet(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	// Get applet from database
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{AccountUUID: account.UUID, UUID: appletId}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applet found",
		})
	}

	// Stop trigger if running
	store.StopTrigger(appletId)
	store.RemoveTrigger(appletId)

	// Delete applet from database
	if result := postgres.DB.Delete(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet deleted",
		},
	})
}

// METHOD: PUT
func StartApplet(c *fiber.Ctx) error {
	account := c.Locals("account").(models.Account)
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	// Get applet from database
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{
		AccountUUID: account.UUID,
		UUID:        appletId,
		Status:      "stopped",
		State:       "complete",
	}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applet found / Applet already started",
		})
	}

	if err := store.StartTrigger(appletId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	if result := postgres.DB.Model(&applet).Updates(&models.Applet{
		Status: "running",
	}); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet started",
		},
	})
}

// It stops an applet by removing the trigger from the database
func StopApplet(c *fiber.Ctx) error {

	account := c.Locals("account").(models.Account)
	appletId, err := uuid.Parse(c.Params("applet_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  fiber.StatusBadRequest,
			"error": "Invalid Applet ID",
		})
	}

	// Get applet from database
	var applet models.Applet
	if result := postgres.DB.Where(&models.Applet{
		AccountUUID: account.UUID,
		UUID:        appletId,
		Status:      "running",
		State:       "complete",
	}).First(&applet); result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"code":  fiber.StatusNoContent,
			"error": "No Applet found / Applet already stopped",
		})
	}

	if err := store.StopTrigger(appletId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"message": "Applet stopped",
		},
	})
}
