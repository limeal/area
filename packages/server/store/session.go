package store

import (
	redis "area-server/db/redis"
	"errors"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

// It's creating a new session store with the given configuration.
var SessionStore = session.New(session.Config{
	Expiration:     3600 * 24 * 7, // 7 days
	KeyLookup:      "cookie:X-Session",
	CookieHTTPOnly: false,
	Storage:        redis.CreateRedisStorage(),
})

// TODO: implement csrf
// Return session ID
func GenerateSession(session *session.Session, accountUUID uuid.UUID) error {
	//crsfToken := utils.GenerateRandomString(32)

	if session.Get("account") != nil {
		return errors.New("Session is invalide")
	}

	session.Set("account", accountUUID.String())

	// Save session
	if err := session.Save(); err != nil {
		return errors.New("Session can't be saved")
	}

	return nil
}

// DestroySession destroys a session
func DestroySession(session *session.Session) bool {
	if err := session.Destroy(); err != nil {
		return false
	}
	return true
}
