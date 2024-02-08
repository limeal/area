package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/store/webhooks"
)

// If the history length has changed, return the ID of the last history item and the name of the author
func hasHistoryBeenUpdated(req static.AreaRequest) shared.AreaResponse {

	historyLen := webhooks.History.Len()

	if (*req.Store)["ctx:history:len"] == nil {
		(*req.Store)["ctx:history:len"] = historyLen
		return shared.AreaResponse{Success: false}
	}

	if (*req.Store)["ctx:history:len"].(int) == historyLen {
		return shared.AreaResponse{Success: false}
	}

	(*req.Store)["ctx:history:len"] = historyLen
	lastHistory := webhooks.History.Front().Value.(webhooks.HistoryItem)
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"webhook:history:id":  lastHistory.ID,
			"webhook:author:name": lastHistory.Author,
		},
	}
}

// `DescriptorForWebhookActionHistoryUpdated` returns a `static.ServiceArea` that describes a webhook
// action history updated event
func DescriptorForWebhookActionHistoryUpdated() static.ServiceArea {
	return static.ServiceArea{
		Name:        "history_updated",
		Description: "Triggered when a webhook action history is updated",
		Method:      hasHistoryBeenUpdated,
		Components: []string{
			"webhook:history:id",
			"webhook:author:name",
		},
	}
}
