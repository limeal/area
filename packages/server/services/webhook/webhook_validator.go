package webhook

import (
	"area-server/classes/static"
	"area-server/services/common"
)

// `WebhookValidators` returns a `static.ServiceValidator` that validates the `url` field of a
// `webhook` request
func WebhookValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:webhook:url": common.URLValidator,
	}
}
