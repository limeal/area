package static

import (
	"area-server/classes/shared"
	"area-server/db/postgres/models"

	"github.com/google/uuid"
)

// `AreaRequest` is a struct that contains the following fields: `AppletID`, `Authorization`,
// `Service`, `Logger`, `Store`, `AuthStore`, and `ExternalData`.
// @property AppletID - The ID of the applet that is being executed.
// @property Authorization - The authorization object that was created when the user logged in.
// @property Service - The service that is being requested.
// @property Logger - A logger that can be used to log messages.
// @property Store - This is the store that is passed to the applet. It is a map of string to
// interface{}.
// @property AuthStore - This is a map of key/value pairs that are stored in the authorization.
// @property ExternalData - This is the data that is passed to the applet from the outside world.
type AreaRequest struct {
	AppletID      uuid.UUID
	Authorization *models.Authorization
	Service       *Service
	Logger        *shared.Logger
	Store         *map[string]interface{}
	AuthStore     map[string]interface{}
	ExternalData  map[string]interface{}
}
