package gmail

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/common"
)

// It returns a map of validators for the Gmail service
func GmailValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:user:id":           UserIDValidator,
		"req:mail:from":         common.EmailValidator,
		"req:mail:to":           common.EmailValidator,
		"req:mail:include:spam": common.BoolValidator,
		"req:auto:reply":        common.BoolValidator,
	}
}

// ------------------------------ Validators ------------------------------

// It takes an authorization, a service, a value, and a store, and returns true if the value is a valid
// user ID
func UserIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{gmail:user:id}}" {
		return true
	}

	_, _, err := service.Endpoints["GmailGetProfileEndpoint"].CallEncode([]interface{}{
		authorization,
		value,
	})

	if err != nil {
		return false
	}

	return true
}
