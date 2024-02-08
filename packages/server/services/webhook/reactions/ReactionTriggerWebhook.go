package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"errors"
	"strconv"
)

// It makes a request to a webhook url with a body and headers
func triggerWebhook(req static.AreaRequest) shared.AreaResponse {

	// Convert body to json string
	url := utils.GenerateFinalComponent((*req.Store)["req:webhook:url"].(string), req.ExternalData, []string{
		".+url",
	})
	body := ""
	if (*req.Store)["req:webhook:body"] != nil {
		body = utils.GenerateFinalComponent((*req.Store)["req:webhook:body"].(string), req.ExternalData, []string{})
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	if (*req.Store)["req:webhook:content:type"] != nil {
		headers["Content-Type"] = (*req.Store)["req:webhook:content:type"].(string)
	}

	method := "POST"
	if (*req.Store)["req:webhook:method"] != nil {
		method = (*req.Store)["req:webhook:method"].(string)
	}

	resp, errr := utils.MakeRequest(url, &utils.RequestParams{
		Method:  method,
		Headers: headers,
		Body:    string(body),
	})

	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	if (*req.Store)["req:webhook:response:status"] != nil {
		status, err := strconv.Atoi((*req.Store)["req:webhook:response:status"].(string))
		if err != nil {
			return shared.AreaResponse{Error: errors.New("Webhook response status is not a number")}
		}
		if resp.StatusCode != status {
			return shared.AreaResponse{Success: false}
		}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea struct that describes the webhook trigger service area
func DescriptorForWebhookReactionTriggerWebhook() static.ServiceArea {
	return static.ServiceArea{
		Name:        "trigger_webhook",
		Description: "Trigger a webhook",
		RequestStore: map[string]static.StoreElement{
			"req:webhook:url": {
				Type:        "string",
				Description: "The URL of the webhook to trigger",
				Required:    true,
			},
			"req:webhook:response:status": {
				Priority:    4,
				Type:        "select",
				Description: "The status of the webhook to match (default: 200)",
				Required:    false,
			},
			"req:webhook:body": {
				Priority:    3,
				Type:        "long_string",
				Description: "The body of the webhook to trigger (default: empty)",
				Required:    false,
			},
			"req:webhook:method": {
				Priority:    1,
				Type:        "select",
				Description: "The method of the webhook to trigger (default: POST)",
				Required:    false,
				Values: []string{
					"POST",
					"GET",
					"PUT",
					"DELETE",
				},
			},
			"req:webhook:content:type": {
				Priority:    2,
				Type:        "select",
				Description: "The content type of the webhook to trigger (default: application/json)",
				Required:    false,
				Values: []string{
					"application/json",
					"application/x-www-form-urlencoded",
					"multipart/form-data",
					"text/plain",
				},
			},
		},
		Method: triggerWebhook,
	}
}
