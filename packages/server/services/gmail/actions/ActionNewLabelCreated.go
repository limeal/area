package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// `ListLabelsResponse` is a struct with a field `Labels` which is a slice of structs with fields `ID`,
// `Name`, `Type`, `MessagesTotal`, `MessagesUnread`, `ThreadsTotal`, `ThreadsUnread`, and `Color`
// which is a struct with fields `Background` and `Text`.
// @property {[]struct {
// 		ID             string `json:"id"`
// 		Name           string `json:"name"`
// 		Type           string `json:"type"`
// 		MessagesTotal  int    `json:"messagesTotal"`
// 		MessagesUnread int    `json:"messagesUnread"`
// 		ThreadsTotal   int    `json:"threadsTotal"`
// 		ThreadsUnread  int    `json:"threadsUnread"`
// 		Color          struct {
// 			Background string `json:"backgroundColor"`
// 			Text       string `json:"textColor"`
// 		} `json:"color"`
// 	}} Labels - An array of label objects.
type ListLabelsResponse struct {
	Labels []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		MessagesTotal  int    `json:"messagesTotal"`
		MessagesUnread int    `json:"messagesUnread"`
		ThreadsTotal   int    `json:"threadsTotal"`
		ThreadsUnread  int    `json:"threadsUnread"`
		Color          struct {
			Background string `json:"backgroundColor"`
			Text       string `json:"textColor"`
		} `json:"color"`
	} `json:"labels"`
}

// It gets the list of labels, checks if the number of labels is the same as the one stored in the
// store, and if it is, it returns the new label
func onNewLabelCreated(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["req:user:id"] != nil {
		(*req.Store)["ctx:user:id"] = (*req.Store)["req:user:id"]
	} else {
		(*req.Store)["ctx:user:id"] = "me"
	}

	encode, _, err := req.Service.Endpoints["ListLabelsEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var response ListLabelsResponse
	if err := json.Unmarshal(encode, &response); err != nil {
		return shared.AreaResponse{Error: err}
	}

	nbLabels := len(response.Labels)
	if ok, errr := utils.IsLatestBasic(req.Store, nbLabels); errr != nil || !ok {
		return shared.AreaResponse{Error: errr, Success: false}
	}

	newLabel := response.Labels[nbLabels-1]
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"gmail:user:id":                (*req.Store)["ctx:user:id"],
			"gmail:label:id":               newLabel.ID,
			"gmail:label:name":             newLabel.Name,
			"gmail:label:type":             newLabel.Type,
			"gmail:label:messages:total":   newLabel.MessagesTotal,
			"gmail:label:messages:unread":  newLabel.MessagesUnread,
			"gmail:label:threads:total":    newLabel.ThreadsTotal,
			"gmail:label:threads:unread":   newLabel.ThreadsUnread,
			"gmail:label:color:background": newLabel.Color.Background,
			"gmail:label:color:text":       newLabel.Color.Text,
		},
	}
}

// `DescriptorForGmailActionNewLabelCreated` returns a `static.ServiceArea` that describes the
// `onNewLabelCreated` function
func DescriptorForGmailActionNewLabelCreated() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_label_created",
		Description: "When a new label is created",
		RequestStore: map[string]static.StoreElement{
			"req:user:id": {
				Description: "The user id to check for new labels",
				Type:        "string",
				Required:    false,
			},
		},
		Method: onNewLabelCreated,
		Components: []string{
			"gmail:label:id",
			"gmail:label:name",
			"gmail:label:type",
			"gmail:label:messages:total",
			"gmail:label:messages:unread",
			"gmail:label:threads:total",
			"gmail:label:threads:unread",
			"gmail:label:color:background",
			"gmail:label:color:text",
		},
	}
}
