package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

/*
 *
 * Example of Authorization:
 *
 * Spotify:
 * UUID: <uuid>
 * AccountUUID: <account_uuid>
 * Service: spotify
 * AccessToken: <access_token>
 * RefreshToken: null
 * ExpiresAt: 2021-08-01 00:00:00 +0000 UTC
 *
 *
 * Area (local):
 * UUID: <uuid>
 * AccountUUID: <account_uuid>
 * Service: @local
 * AccessToken: <jwt token containing email and hashed password of the user>
 * RefreshToken: null
 * ExpiresIn: -1
 *
 * Twitch:
 * UUID: <uuid>
 * AccountUUID: <account_uuid>
 * Service: twitch
 * AccessToken: <access_token>
 * RefreshToken: <refresh_token>
 * ExpiresAt: 2021-08-01 00:00:00 +0000 UTC
 *
 *
 * Permanent option is used to know if the authorization is linked directly to the account or not.
 * Example: if the user register itself with the Service, the authorization is permanent.
 * but if the user link the service to his account, the authorization is not permanent.
 */

// Authorization -> Many to One -> Account

type Authorization struct {
	UUID         uuid.UUID      `gorm:"primaryKey" json:"-"`
	Account      Account        `gorm:"foreignKey:AccountUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AccountUUID  uuid.UUID      `gorm:"not null" json:"-"`
	Type         string         `gorm:"not null" json:"type"` // "oauth2", ...
	AuthService  string         `gorm:"not null" json:"name"` // authenticator
	AccessToken  string         `gorm:"not null" json:"-"`
	RefreshToken string         `json:"-"`
	Other        datatypes.JSON `json:"-"`
	Permanent    bool           `gorm:"default:false" json:"permanent"` // Permanent=true, means that it can't be deleted
	ExpireAt     time.Time      `gorm:"not null" json:"expire_at"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"-"`
}
