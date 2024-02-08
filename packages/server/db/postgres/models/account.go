package models

import (
	"time"

	"github.com/google/uuid"
)

// Account is a struct that has a UUID, Authenticator, Email, Username, Password, CreatedAt, and
// UpdatedAt field.
// @property UUID - The UUID of the account.
// @property {string} Authenticator - The service that the account is associated with.
// @property {string} Email - The email address of the account.
// @property {string} Username - The username of the account.
// @property {string} Password - The password for the account. This is only used for local accounts.
// @property CreatedAt - The time the account was created.
// @property UpdatedAt - This is the time the account was last updated.
type Account struct {
	UUID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Authenticator string    `gorm:"not null" json:"-"`
	Email         string    `gorm:"unique" json:"email"`
	Username      string    `gorm:"default:'noob'" json:"username"`
	Password      string    `json:"-"` // Only used for local accounts (Service = "@local")
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"-"`
}
