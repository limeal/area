package triggers

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/services"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
	Area         *context.TriggerArea
	Logger       *Logger
	ExternalData map[string]interface{}
*/

// `TriggerArea` is a struct that contains a `Model` of type `*models.Area`, an `Authorization` of type
// `*models.Authorization`, a `Service` of type `*static.Service`, an `Area` of type
// `*static.ServiceArea`, a `Store` of type `map[string]interface{}`, an `AuthStore` of type
// `map[string]interface{}`, and a `SnapshotID` of type `uuid.UUID`.
// @property Model - The Area model from the database.
// @property Authorization - The authorization that triggered the area.
// @property Service - The service name
// @property Area - The area that the trigger is in.
// @property Store - This is a map of key/value pairs that can be used to store data for the trigger.
// @property AuthStore - This is a map of the authorization store.
// @property SnapshotID - This is the ID of the snapshot that was used to create this trigger area. If
// the snapshot changes, this ID will change.
type TriggerArea struct {
	Model         *models.Area
	Authorization *models.Authorization
	Service       *static.Service     // Service name
	Area          *static.ServiceArea // Area name
	Store         map[string]interface{}
	AuthStore     map[string]interface{}
	SnapshotID    uuid.UUID // Snapshot ID (If changed, call Update())
}

// `TriggerResponse` is a struct with three fields: `Error`, `Success`, and `Data`.
//
// The `Error` field is of type `error`, which is a built-in Go type.
//
// The `Success` field is of type `bool`, which is a built-in Go type.
//
// The `Data` field is of type `map[string]interface{}`, which is a built-in Go type.
// @property {error} Error - This is the error that occurred during the execution of the trigger.
// @property {bool} Success - This is a boolean value that indicates whether the trigger was successful
// or not.
// @property Data - This is the data that will be passed to the trigger.
type TriggerResponse struct {
	Error   error
	Success bool
	Data    map[string]interface{}
}

// > This function creates a new trigger area
func NewArea(appletID uuid.UUID, area *models.Area, areatype string) (*TriggerArea, error) {
	areaModel := &TriggerArea{}
	if err := areaModel.Update(appletID, area, areatype); err != nil {
		return nil, err
	}
	return areaModel, nil
}

// This function is updating the trigger area.
func (a *TriggerArea) Update(appletID uuid.UUID, area *models.Area, areatype string) error {
	// 1. Get Authorization

	if a.Service = services.GetServiceByName(area.Service); a.Service == nil {
		return errors.New("Area: Service not found :>" + area.Service)
	}

	a.Model = area

	var authorization models.Authorization
	if a.Service.Authenticator != nil {
		if result := postgres.DB.Where(&models.Authorization{
			UUID: area.AuthorizationUUID,
		}).First(&authorization); result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
			return errors.New("Area: Authorization not found !")
		}
		a.Authorization = &authorization

		// 3. Get AuthStore
		if a.Authorization.Other != nil {
			authStore := make(map[string]interface{})

			if err := json.Unmarshal([]byte(a.Authorization.Other.String()), &authStore); err != nil {
				return errors.New("Area: Auth Store is not valid !")
			}
			a.AuthStore = authStore
		}

	} else {
		a.Authorization = nil
	}

	// 3. Get Store
	store := make(map[string]interface{})

	if err := json.Unmarshal([]byte(a.Model.Store.String()), &store); err != nil {
		return errors.New("Area: Store is not valid !")
	}

	var elem *static.ServiceArea
	// 4. Get Action or Reaction
	if a.Model.Type == "action" {
		if elem = a.Service.GetActionByName(a.Model.Name); elem == nil {
			return errors.New("Area: Action not found :>" + a.Model.Name)
		}
	} else {
		if elem = a.Service.GetReactionByName(a.Model.Name); elem == nil {
			return errors.New("Area: Reaction not found :>" + a.Model.Name)
		}
	}
	a.Area = elem

	a.Store = static.MergeStore(elem.RequestStore, store)

	// 5. Generate Snapshot ID
	a.SnapshotID = uuid.New()
	return nil
}

// Calling the `Method` function of the `Area` struct.
func (a *TriggerArea) Call(
	appletID uuid.UUID,
	logger *shared.Logger,
	externaldata map[string]interface{},
) shared.AreaResponse {
	return a.Area.Method(static.AreaRequest{
		AppletID:      appletID,
		Authorization: a.Authorization,
		Service:       a.Service,
		Logger:        logger,
		Store:         &a.Store,
		AuthStore:     a.AuthStore,
		ExternalData:  externaldata,
	})
}

// This function is refreshing the authorization.
func (a *TriggerArea) Refresh() error {
	if a.Service.Authenticator == nil {
		return nil
	}
	var err error
	a.Authorization, err = a.Service.Authenticator.RefreshToken(a.Authorization)
	return err
}
