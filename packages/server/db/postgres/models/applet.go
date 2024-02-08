package models

import (
	"time"

	"github.com/google/uuid"
)

/*
 * Example of an applet:
 *
 * An app that post a message on a discord channel when a track is added to playlist:
 * UUID: <uuid>
 * AccountUUID: <account_uuid>
 * Name: "Discord notification"
 * Description: "Send a message on a discord channel when a track is added to playlist"
 * Emitter: "service:spotify;name:track_added_to_playlist" - Must be an Action
 * EmitterSettings: "query_1:<playlist_id>" - Must be an Action
 * Receiver: "service:discord;name:send_message" - Must be a Reaction
 * ReceiverSettings: "query_1:<channel_id>;query_2:<message>" - Must be a Reaction
 * Public: false - When the applet is public, it can be used by anyone and displayed on the applet store
 */

// Applet -> Many to One -> Account
type Applet struct {
	UUID        uuid.UUID `gorm:"primaryKey" json:"id"`                                                                         // Unique UUID of the applet
	Account     Account   `gorm:"foreignKey:AccountUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"` // Many to One (Account linked to the applet)
	AccountUUID uuid.UUID `gorm:"not null" json:"account_uuid"`                                                                 // Unique
	Action      string    `gorm:"default:null" json:"action"`                                                                   // <service>;<name> of applet action
	Name        string    `json:"name"`                                                                                         // Name of the applet
	Description string    `json:"description"`                                                                                  // Description of the applet
	State       string    `gorm:"not null" json:"-"`                                                                            // State of the applet (partial, complete)
	Public      bool      `gorm:"default:false" json:"public"`                                                                  // If the applet is public, it can be copied by anyone
	Active      bool      `gorm:"default:true" json:"active"`                                                                   // If the applet is active, it can be triggered
	Status      string    `gorm:"default:'stopped'" json:"status"`                                                              // Status of the applet (stopped, running)
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
}
