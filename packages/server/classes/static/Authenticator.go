package static

import (
	"area-server/db/postgres"
	"area-server/db/postgres/models"
	"area-server/utils"
	"errors"
	"fmt"
	"time"
)

// Authenticator is an interface that has three methods: Authenticate, ReAuthenticate, and Disprove.
// @property {error} Authenticate - This is the method that will be called when the user is trying to
// authenticate.
// @property {error} ReAuthenticate - This is used to re-authenticate the user. This is used when the
// user is already authenticated but the session has expired.
// @property {error} Disprove - This is the method that will be called when the user is trying to
// disprove the authentication.
type Authenticator interface {
	Authenticate([]interface{}) error
	ReAuthenticate() error
	Disprove() error
}

// `AuthEndpoints` is a struct that contains a `RequestDescriptor` named `AccessToken` and four
// optional `RequestDescriptor`s named `RefreshToken`, `RevokeToken`, `ValidateToken`, and `Profile`.
// @property AccessToken - The endpoint that will be used to request an access token.
// @property RefreshToken - The endpoint to refresh the access token.
// @property RevokeToken - This is the endpoint that will be called to revoke the access token.
// @property ValidateToken - This is the endpoint that will be called to validate the token.
// @property Profile - The endpoint to get the user's profile.
// @property Email - The endpoint to get the user's email address.
type AuthEndpoints struct {
	AccessToken   utils.RequestDescriptor  `json:"access_token"`
	RefreshToken  *utils.RequestDescriptor `json:"refresh_token,omitempty"`
	RevokeToken   *utils.RequestDescriptor `json:"revoke_token,omitempty"`
	ValidateToken *utils.RequestDescriptor `json:"validate_token,omitempty"`
	Profile       *utils.RequestDescriptor `json:"profile,omitempty"`
	Email         *utils.RequestDescriptor `json:"email,omitempty"`
}

// `OAuth2Authenticator` is a struct with a bunch of fields.
// @property {string} Name - The name of the authenticator.
// @property {bool} Enabled - If true, the service can be used for authentication
// @property {More} More - This is a struct that contains the following properties:
// @property {string} AuthorizationURI - The URI to which the user is redirected to authorize the
// application.
// @property {AuthEndpoints} AuthEndpoints - This is a struct that contains the endpoints for the
// OAuth2 flow.
// @property OtherParams - This is a function that takes a map of string to interface and returns a map
// of string to interface. This is used to add additional parameters to the OAuth2 request.
type OAuth2Authenticator struct {
	Name             string                                              `json:"name"`
	Enabled          bool                                                `json:"enabled"` // If true, the service can be used for authentication
	More             More                                                `json:"more"`
	AuthorizationURI string                                              `json:"authorization_uri"`
	AuthEndpoints    AuthEndpoints                                       `json:"-"`
	OtherParams      func(map[string]interface{}) map[string]interface{} `json:"-"`
}

// AuthenticateError is a struct with two fields, StatusCode and ErrorDesc.
// @property {int} StatusCode - The HTTP status code returned by the server.
// @property {string} ErrorDesc - This is the error description that will be returned to the user.
type AuthenticateError struct {
	StatusCode int
	ErrorDesc  string
}

// AuthenticateResponse is a struct with two fields, Type and Data.
//
// The Type field is a string, and the Data field is a map of strings to interfaces.
//
// The interface{} type is a special type in Go that can hold any value.
// @property {string} Type - The type of authentication.
// @property Data - This is a map of key-value pairs that are specific to the type of authentication.
type AuthenticateResponse struct {
	Type string `json:"type"` // "oauth2", "basic", "custom"
	Data map[string]interface{}
}

// A function that takes a slice of interfaces and returns a pointer to an AuthenticateResponse and a
// pointer to an AuthenticateError.
func (a *OAuth2Authenticator) Authenticate(params []interface{}) (*AuthenticateResponse, *AuthenticateError) {

	needAuth, okO := params[0].(bool)
	if !okO {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "needAuth is not a boolean"}
	}
	code, okT := params[1].(string)
	if !okT {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "code is not a string"}
	}
	redirectURI, okTh := params[2].(string)
	if !okTh {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "redirect_uri is not a string"}
	}

	if needAuth && a.Enabled == false {
		return nil, &AuthenticateError{StatusCode: 400, ErrorDesc: "Service does not provide authentication"}
	}

	if needAuth && a.Enabled == true && a.AuthEndpoints.Email == nil {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "Service email endpoint not found"}
	}

	authr, _, err := a.AuthEndpoints.AccessToken.Call([]interface{}{code, redirectURI})

	if err != nil {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "Service could not fetch access token"}
	}

	oauth := make(map[string]interface{})
	ok, oke := false, false
	oauth["access_token"], ok = authr["access_token"]
	oauth["token_type"], oke = authr["token_type"]

	if !ok || !oke {
		return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: errors.New("access_token or token_type not found").Error()}
	}

	oauth["refresh_token"] = authr["refresh_token"]
	oauth["expires_in"] = authr["expires_in"]

	email := ""
	if needAuth {
		content, _, err := a.AuthEndpoints.Email.Call([]interface{}{oauth})

		if err != nil {
			return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: err.Error()}
		}

		if content["email"] == nil {
			return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: "email not found"}
		}

		email, _ = content["email"].(string)
	}

	if a.AuthEndpoints.ValidateToken != nil {
		content, _, err := a.AuthEndpoints.ValidateToken.Call([]interface{}{oauth})

		if err != nil {
			return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: err.Error()}
		}

		authr = utils.MergeMaps(authr, content)
	}

	if a.AuthEndpoints.Profile != nil {
		content, _, err := a.AuthEndpoints.Profile.Call([]interface{}{oauth})

		if err != nil {
			return nil, &AuthenticateError{StatusCode: 500, ErrorDesc: err.Error()}
		}

		authr = utils.MergeMaps(authr, content)
	}

	var expiredAt time.Time
	if oauth["expires_in"] == nil {
		expiredAt = time.Unix(4072721567, 0)
	} else {
		expiredAt = time.Now().Add(time.Duration(int(oauth["expires_in"].(float64))) * time.Second)
	}

	var refreshToken string
	if oauth["refresh_token"] == nil {
		refreshToken = ""
	} else {
		refreshToken = oauth["refresh_token"].(string)
	}

	var other map[string]interface{}
	if a.OtherParams != nil {
		other = a.OtherParams(authr)
	}

	return &AuthenticateResponse{
		Type: "oauth2",
		Data: map[string]interface{}{
			"access_token":  oauth["access_token"],
			"token_type":    oauth["token_type"],
			"refresh_token": refreshToken,
			"expired_at":    expiredAt,
			"email":         email,
			"other":         other,
		},
	}, nil
}

// Refresh Token
func (a *OAuth2Authenticator) RefreshToken(authorization *models.Authorization) (*models.Authorization, error) {
	if authorization == nil {
		return nil, fmt.Errorf("Service: Authorization is nil")
	}

	if !authorization.ExpireAt.Before(time.Now()) {
		return authorization, nil
	}

	if a.AuthEndpoints.RefreshToken == nil {
		return nil, fmt.Errorf("Service: Refresh token is not supported")
	}

	response, _, err := a.AuthEndpoints.RefreshToken.Call([]interface{}{authorization.RefreshToken})
	if err != nil {
		return nil, err
	}

	if result := postgres.DB.Model(authorization).Update("access_token", response["access_token"].(string)); result.Error != nil {
		return nil, result.Error
	}

	authorization.AccessToken = response["access_token"].(string)

	// Update refresh token if provided
	if response["refresh_token"] != nil {
		authorization.RefreshToken = response["refresh_token"].(string)
		if result := postgres.DB.Model(authorization).Update("refresh_token", response["refresh_token"].(string)); result.Error != nil {
			return nil, result.Error
		}
	}

	authorization.ExpireAt = time.Now().Add(time.Duration(response["expires_in"].(float64)) * time.Second)
	return authorization, nil
}
