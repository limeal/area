package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/json"
)

// `SubscribeSnippet` is a struct that contains a `ResourceId` field, which is a struct that contains a
// `ChannelId` field, which is a string.
// @property ResourceId - The channel ID of the channel that the user subscribed to.
type SubscribeSnippet struct {
	ResourceId struct {
		ChannelId string `json:"channelId"`
	} `json:"resourceId"`
}

// It adds a subscription to the channel with the ID stored in the `req:channel:id` key of the
// request's store
func addSubscriptionToChannel(req static.AreaRequest) shared.AreaResponse {

	channelID := (*req.Store)["req:channel:id"].(string)

	ebody, err := json.Marshal(map[string]interface{}{
		"snippet": SubscribeSnippet{
			ResourceId: struct {
				ChannelId string `json:"channelId"`
			}{
				ChannelId: channelID,
			},
		},
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, _, err = req.Service.Endpoints["AddSubscriptionEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		string(ebody),
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a static.ServiceArea struct that describes the service area
func DescriptorForYoutubeReactionAddSubscription() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_subscription",
		Description: "Add a subscription to a channel",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "string",
				Description: "Name of the channel",
				Required:    true,
			},
		},
		Method: addSubscriptionToChannel,
	}
}
