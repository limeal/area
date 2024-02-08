package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

/*
 * Example of an Action / Reaction:
 *
 * An app that post a message on a discord channel when a track is added to playlist:
 * UUID: <uuid>
 * AppletUUID: <applet_uuid>
 * AuthorizationUUID: <authorization_uuid>
 * Name: "Discord notification"
 * Store: { "query_1": "<channel_id>", "query_2": "<message>" }
 */

// Area -> One to One -> Applet
type Area struct {
	UUID              uuid.UUID      `gorm:"primaryKey" json:"id"`
	Applet            Applet         `gorm:"foreignKey:AppletUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AppletUUID        uuid.UUID      `gorm:"not null" json:"-"`
	AuthorizationUUID uuid.UUID      `gorm:"not null" json:"-"`
	Type              string         `gorm:"not null" json:"type"`
	Service           string         `gorm:"not null" json:"service"`
	Name              string         `gorm:"not null" json:"name"`
	Store             datatypes.JSON `gorm:"type:jsonb;not null" json:"store"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"-"`
}
