package static

import (
	"area-server/db/postgres/models"
	"area-server/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*
 * Example:
 * Name: "Spotify"
 * BaseURI: "https://api.spotify.com/v1"
 * Endpoints:
 *  - "authorize": "/authorize"
 *  - "me": "/me"
 */

// `ServiceRoute` is a struct that contains a string, a function, a string, and a bool.
// @property {string} Endpoint - The endpoint that the route will be registered to.
// @property Handler - The function that will be called when the route is hit.
// @property {string} Method - The HTTP method to use for this route.
// @property {bool} NeedAuth - If true, the request will be checked for a valid JWT token.
type ServiceRoute struct {
	Endpoint string                 `json:"endpoint"`
	Handler  func(*fiber.Ctx) error `json:"-"`
	Method   string                 `json:"method"`
	NeedAuth bool                   `json:"need_auth"`
}

// `ServiceValidator` is a map of strings to functions that take an `Authorization`, a `Service`, an
// interface, and a map of strings to interfaces and return a boolean.
//
type ServiceValidator (map[string]func(*models.Authorization, *Service, interface{}, map[string]interface{}) bool)
type ServiceEndpoint (map[string]*utils.RequestDescriptor)

// More is a struct with two fields, Avatar and Color, both of which are exported.
// @property {bool} Avatar - If true, the bot will take the avatar of the user who sent the message.
// @property {string} Color - The color of the text.
type More struct {
	Avatar bool   `json:"take_avatar"`
	Color  string `json:"color"`
}

type Service struct {
	Name          string               `json:"name"` // Name of the service
	Description   string               `json:"description"`
	Authenticator *OAuth2Authenticator `json:"authenticator"` // Authenticator used by the service
	RateLimit     float64              `json:"-"`             // Number of requests per 30 seconds
	More          *More                `json:"more"`          // If equal to nil, the client must use the avatar_uri and color of the authenticator (if exists)
	Validators    ServiceValidator     `json:"-"`             // Validators used by the service
	Endpoints     ServiceEndpoint      `json:"-"`             // Endpoints used by the service
	Routes        []ServiceRoute       `json:"-"`             // Routes used by the service
	Gateway       Gateway              `json:"-"`             // Gateway used by the service
	Actions       []ServiceArea        `json:"actions"`       // Actions provided by the service
	Reactions     []ServiceArea        `json:"reactions"`     // Reactions provided by the service
}

// It's a method that returns a string representation of the service.
func (s *Service) String() string {
	return fmt.Sprint(fiber.Map{
		"name":          s.Name,
		"description":   s.Description,
		"authenticator": s.Authenticator,
		"rate_limit":    s.RateLimit,
		"more":          s.More,
		"actions":       s.Actions,
		"reactions":     s.Reactions,
	})
}

// It's a method that returns a pointer to a `ServiceArea` struct.
func (s *Service) GetActionByName(name string) *ServiceArea {
	for _, action := range s.Actions {
		if action.Name == name {
			return &action
		}
	}
	return nil
}

// It's a method that returns a pointer to a `ServiceArea` struct.
func (s *Service) GetReactionByName(name string) *ServiceArea {
	for _, reaction := range s.Reactions {
		if reaction.Name == name {
			return &reaction
		}
	}
	return nil
}

// It's a method that returns a slice of `ServiceArea` structs.
func (s *Service) GetActions() []ServiceArea {
	return s.Actions
}

// It's returning a slice of `ServiceArea` structs.
func (s *Service) GetReactions() []ServiceArea {
	return s.Reactions
}

// RefreshTokenResponse is a type that represents a response from the refresh token endpoint.
// @property {string} AccessToken - The new access token.
// @property ExpireAt - The time when the access token will expire.
type RefreshTokenResponse struct {
	AccessToken string    `json:"access"`
	ExpireAt    time.Time `json:"expire_at"`
}
