package main

import (
	"area-server/classes/triggers"
	"area-server/store"
	"area-server/store/webhooks"

	"area-server/db/postgres"
	"area-server/db/postgres/models"
)

// It loads all the applets from the database and adds them to the store
func LoadApplets() error {
	var applets []models.Applet
	if result := postgres.DB.Where(&models.Applet{State: "complete"}).Find(&applets); result.Error != nil {
		if result.RowsAffected == 0 {
			return nil
		}
		return result.Error
	}

	for _, applet := range applets {

		var action models.Area
		if result := postgres.DB.Where(&models.Area{AppletUUID: applet.UUID, Type: "action"}).First(&action); result.Error != nil {
			return result.Error
		}
		if action.Service == "webhook" && action.Name == "applet_triggered" {
			webhooks.AddWebhook(applet.UUID.String())
		}

		tr, err := triggers.CreateTrigger(&applet)
		if err != nil {
			return err
		}

		store.AddTrigger(tr, applet.Status != "stopped")
	}

	return nil
}
