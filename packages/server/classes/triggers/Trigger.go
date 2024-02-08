package triggers

import (
	"area-server/classes/shared"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// A Trigger is a struct that contains an AppletID, an EmitterArea, a slice of ReceiversArea, a boolean
// Active, a boolean Stopped, a channel of bool OnOff, a channel of bool Interrupt, and a channel of
// string slice Update.
// @property AppletID - The ID of the applet that owns this trigger.
// @property EmitterArea - The area that the trigger is emitting from.
// @property {[]*TriggerArea} ReceiversArea - This is a slice of TriggerArea structs. This is the area
// where the trigger will be active.
// @property {bool} Active - This is a boolean value that indicates whether the trigger is active or
// not.
// @property {bool} Stopped - This is a boolean value that indicates whether the trigger is stopped or
// not.
// @property OnOff - This channel is used to turn the trigger on and off.
// @property Interrupt - This channel is used to interrupt the trigger.
// @property Update - This channel is used to update the trigger's receivers.
type Trigger struct {
	AppletID      uuid.UUID
	EmitterArea   *TriggerArea
	ReceiversArea []*TriggerArea
	Active        bool
	Stopped       bool
	OnOff         chan bool
	Interrupt     chan bool
	Update        chan []string
}

// It creates a new trigger for an applet
func CreateTrigger(app *models.Applet) (*Trigger, error) {
	// 1. Get Emitter Area

	var action models.Area
	if result := postgres.DB.Where(&models.Area{AppletUUID: app.UUID, Type: "action"}).First(&action); result.Error != nil {
		return nil, result.Error
	}

	emitter, err := NewArea(app.UUID, &action, "action")
	if err != nil {
		return nil, err
	}

	var reactions []models.Area
	if result := postgres.DB.Where(&models.Area{AppletUUID: app.UUID, Type: "reaction"}).Find(&reactions); result.Error != nil {
		return nil, result.Error
	}

	var receivers []*TriggerArea
	// 2. Get Receiver Area
	for _, reaction := range reactions {
		reactionArea, err := NewArea(app.UUID, &reaction, "reaction")
		if err != nil {
			return nil, err
		}
		receivers = append(receivers, reactionArea)
	}

	return &Trigger{
		AppletID:      app.UUID,
		EmitterArea:   emitter,
		ReceiversArea: receivers,
		Active:        true,
		Stopped:       false,
		OnOff:         make(chan bool),
		Interrupt:     make(chan bool),
		Update:        make(chan []string),
	}, nil
}

// Stopping the trigger.
func (t *Trigger) Stop() error {

	if t.EmitterArea.Area.UseGateway && t.EmitterArea.Service.Gateway != nil {
		t.EmitterArea.Service.Gateway.Stop()
	}

	if result := postgres.DB.Model(&models.Applet{UUID: t.AppletID}).Updates(&models.Applet{
		Status: "stopped",
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

// A function that listens to the trigger.
func (t *Trigger) Listen() error {

	logger := shared.NewLogger(t.AppletID)
	if logger == nil {
		return fmt.Errorf("Trigger: Error while creatin logger !")
	}
	defer logger.Close()

	wtime := time.Duration(0)
	logger.WriteInfo("Start Application !", true)

	if t.EmitterArea.Area.UseGateway && t.EmitterArea.Service.Gateway != nil {
		go t.EmitterArea.Service.Gateway.Start()
	}

	for {
		select {
		case <-t.Interrupt:
			logger.WriteInfo("Stop Application !", true)
			return t.Stop() /*
				case data := <-t.Update:
					// TODO: Remake it work
					if err := t.EmitterArea.Update(t.AppletID, data[0], data[1]); err != nil {
						logger.WriteError("Updating " + data[1] + " Failed :> " + err.Error())
						return t.Stop()
					}
					logger.WriteInfo(data[0]+" Updated :> ("+t.EmitterArea.Model.Service+":"+t.EmitterArea.Area.Name+")", true) */
		case data := <-t.OnOff:
			logger.WriteInfo("Application is now "+utils.TernaryOperator(data, "active", "inactive").(string)+" !", true)
			t.Active = data
		case <-time.After(wtime):
			if wtime == 0 && !t.EmitterArea.Area.UseGateway {
				wtime = (time.Duration(30/t.EmitterArea.Service.RateLimit) * time.Second)
			}

			if !t.Active {
				logger.WriteInfo("Application is not active !", false)
				continue
			}

			if err := t.EmitterArea.Refresh(); err != nil {
				logger.WriteError("Refreshing Token Failed !")
				return t.Stop()
			}

			logger.WriteInfo("Check if action is triggered !", false)
			emResponse := t.EmitterArea.Call(t.AppletID, logger, nil)
			if emResponse.Error != nil {
				logger.WriteError("Action provide an error :> " + emResponse.Error.Error())
				return t.Stop()
			}

			if !emResponse.Success {
				logger.WriteInfo("Action not triggered retry in "+wtime.String(), false)
				continue
			}

			logger.WriteInfo("Action Triggered !", true)
			for _, receiver := range t.ReceiversArea {
				if err := receiver.Refresh(); err != nil {
					logger.WriteError("Refreshing Token Failed :> " + err.Error())
					return t.Stop()
				}

				if rcResponse := receiver.Call(t.AppletID, logger, emResponse.Data); rcResponse.Error != nil {
					logger.WriteError("Reaction provide an error :> " + rcResponse.Error.Error())
					return t.Stop()
				}
				logger.WriteInfo("Reaction Triggered !", true)
			}
		}
	}
}
