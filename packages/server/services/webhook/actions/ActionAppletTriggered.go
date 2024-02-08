package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/store/webhooks"
	"time"
)

// It waits for a webhook to be triggered, and if it is, it returns the data that was sent with the
// webhook
func isAppletTriggered(req static.AreaRequest) shared.AreaResponse {

	webhook, err := webhooks.GetWebhook(req.AppletID.String(), webhooks.WebhookTypeAppletTrigger)
	if err != nil {
		return shared.AreaResponse{
			Success: false,
		}
	}

	select {
	case data := <-webhook.Data:
		return shared.AreaResponse{
			Success: true,
			Data: map[string]interface{}{
				"webhook:data": string(data),
			},
		}
	case <-time.After(2 * time.Second):
		return shared.AreaResponse{
			Success: false,
		}
	}
}

// `DescriptorForWebhookActionAppletTriggered` returns a `static.ServiceArea` that describes the
// webhook action applet triggered service area
func DescriptorForWebhookActionAppletTriggered() static.ServiceArea {
	return static.ServiceArea{
		Name:        "applet_triggered",
		Description: "Triggered when a webhook action applet is triggered",
		Method:      isAppletTriggered,
		Components: []string{
			"webhook:data",
		},
	}
}
