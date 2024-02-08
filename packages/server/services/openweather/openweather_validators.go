package openweather

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
)

// It returns a `static.ServiceValidator` which is a map of `string` to `Validator` functions
func OpenWeatherValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:station:id": StationIdValidator,
	}
}

// ------------------ Validators ------------------

// It checks if the value is a string, and if it is, it checks if the value is a valid station ID
func StationIdValidator(
	authorization *models.Authorization,
	service *static.Service,
	value interface{},
	store map[string]interface{},
) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value == "{{openw:station:id}}" {
		return true
	}

	_, _, err := service.Endpoints["GetStationEndpoint"].CallEncode([]interface{}{
		value.(string),
	})

	return err == nil
}
