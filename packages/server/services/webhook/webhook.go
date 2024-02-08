package webhook

import (
	"area-server/classes/static"
	"area-server/services/webhook/actions"
	"area-server/services/webhook/reactions"
)

// It returns a static.Service object that describes the webhook service
func Descriptor() static.Service {
	return static.Service{
		Name:        "webhook",
		Description: "Webhook service",
		RateLimit:   0,
		More: &static.More{
			Avatar: true,
			Color:  "#7289DA",
		},
		Routes:     WebhookRoute(),
		Validators: WebhookValidators(),
		Actions: []static.ServiceArea{
			actions.DescriptorForWebhookActionAppletTriggered(),
			actions.DescriptorForWebhookActionHistoryUpdated(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForWebhookReactionTriggerWebhook(),
		},
	}
}
