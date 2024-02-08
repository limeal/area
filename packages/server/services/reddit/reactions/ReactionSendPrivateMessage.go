package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
)

// It calls the `ComposePrivateMessageEndpoint` endpoint with the `Authorization` header, the
// `req:user:to` store value, the `req:message:subject` store value, and the `req:message:body` store
// value
func sendPrivateMessage(req static.AreaRequest) shared.AreaResponse {

	_, _, err := req.Service.Endpoints["ComposePrivateMessageEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["req:user:to"],
		(*req.Store)["req:message:subject"],
		(*req.Store)["req:message:body"],
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// `DescriptorForRedditReactionNewPostOnSubReddit` returns a `static.ServiceArea` with the name
// `send_private_message`, a description, a request store, and a method
func DescriptorForRedditReactionNewPostOnSubReddit() static.ServiceArea {
	return static.ServiceArea{
		Name:        "send_private_message",
		Description: "Send a private message to a user",
		RequestStore: map[string]static.StoreElement{
			"req:user:to": {
				Required:    true,
				Description: "The name of the user to send the message to",
				Type:        "string",
			},
			"req:message:subject": {
				Priority:    1,
				Required:    true,
				Description: "The subject of the message",
				Type:        "string",
			},
			"req:message:body": {
				Priority:    2,
				Required:    true,
				Description: "The body of the message",
				Type:        "string",
			},
		},
		Method: sendPrivateMessage,
	}
}
