package store

import (
	"area-server/classes/triggers"
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"fmt"

	"github.com/google/uuid"
)

// Creating a map of uuid.UUID to triggers.Trigger pointers.
var Triggers = make(map[uuid.UUID]*triggers.Trigger)

// It adds a trigger to the Triggers map, and if autoStart is true, it starts listening for the trigger
func AddTrigger(trigger *triggers.Trigger, autoStart bool) {
	fmt.Println("Adding trigger for applet: ", trigger.AppletID)
	Triggers[trigger.AppletID] = trigger

	if autoStart {
		go trigger.Listen()
	}
}

// It takes an applet ID, finds the trigger for that applet, and then sends a message to the trigger's
// OnOff channel
func OpenTrigger(id uuid.UUID) error {
	fmt.Println("Opening trigger for applet: ", id)
	_, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	Triggers[id].OnOff <- true

	if result := postgres.DB.Model(&models.Applet{UUID: id}).Update("active", true); result.Error != nil {
		return result.Error
	}
	return nil
}

// It takes an applet ID, finds the trigger for that applet, and sends a false value to the trigger's
// OnOff channel
func CloseTrigger(id uuid.UUID) error {
	fmt.Println("Closing trigger for applet: ", id)
	_, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	Triggers[id].OnOff <- false

	// Update Applet Status in DB
	if result := postgres.DB.Model(&models.Applet{UUID: id}).Update("active", false); result.Error != nil {
		return result.Error
	}
	return nil
}

// It removes a trigger from the `Triggers` map
func RemoveTrigger(id uuid.UUID) error {
	fmt.Println("Removing trigger for applet: ", id)
	_, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	delete(Triggers, id)
	return nil
}

// It starts a trigger for an applet
func StartTrigger(id uuid.UUID) error {
	fmt.Println("Starting trigger for applet: ", id)
	_, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	Triggers[id].Interrupt = make(chan bool)
	Triggers[id].Update = make(chan []string)
	Triggers[id].OnOff = make(chan bool)
	Triggers[id].Stopped = false
	go Triggers[id].Listen()
	return nil
}

// It closes all the channels for the trigger, and sets the trigger's Stopped flag to true
func StopTrigger(id uuid.UUID) error {
	fmt.Println("Stopping trigger for applet: ", id)
	_, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	select {
	case <-Triggers[id].Interrupt:
		break
	default:
		close(Triggers[id].Interrupt)
	}
	select {
	case <-Triggers[id].Update:
		break
	default:
		close(Triggers[id].Update)
	}
	select {
	case <-Triggers[id].OnOff:
		break
	default:
		close(Triggers[id].OnOff)
	}
	Triggers[id].Stopped = true
	return nil
}

func UpdateTrigger(id uuid.UUID, areaType string) error {
	// TODO: Reimplement this
	/* _, ok := Triggers[id]
	if !ok {
		return fmt.Errorf("Trigger not found")
	}
	if !Triggers[id].Stopped {
		Triggers[id].Update <- []string{areaType, "update"}
		return nil
	}
	if areaType == "action" {
		if err := Triggers[id].EmitterArea.Update(Triggers[id].AppletID, areaType); err != nil {
			return err
		}
	} else {
		if err := Triggers[id].ReceiverArea.Update(Triggers[id].AppletID, areaType); err != nil {
			return err
		}
	} */
	return nil
}

// GetTrigger returns a pointer to a Trigger struct if the triggerID is found in the Triggers map,
// otherwise it returns nil.
func GetTrigger(triggerID uuid.UUID) *triggers.Trigger {
	trigger, ok := Triggers[triggerID]
	if !ok {
		return nil
	}
	return trigger
}
