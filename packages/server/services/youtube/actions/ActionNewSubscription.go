package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"
)

// It checks if the user is subscribed to a new channel
func isSubscribingToNewChannel(req static.AreaRequest) shared.AreaResponse {
	query := make(map[string]string)

	query["part"] = "snippet"
	query["maxResults"] = "1"
	query["mine"] = "true"

	if (*req.Store)["req:creator:id"] != nil {
		query["forChannelId"] = (*req.Store)["req:creator:id"].(string)
	}

	encode, _, err := req.Service.Endpoints["GetAllSubscriptionsEndpoint"].CallEncode([]interface{}{req.Authorization, query})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	subList := common.YoutubeSubscriptionResponse{}
	errr := json.Unmarshal(encode, &subList)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}
	nbSubscriptions := subList.PageInfo.TotalResults
	ok, errL := utils.IsLatestBasic(req.Store, nbSubscriptions)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}

	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"youtube:sub:id":           subList.Items[0].ID,
			"youtube:sub:published_at": subList.Items[0].Snippet.PublishedAt,
			"youtube:sub:title":        subList.Items[0].Snippet.Title,
			"youtube:sub:channel":      subList.Items[0].Snippet.ChannelID,
			"youtube:sub:description":  subList.Items[0].Snippet.Description,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area "new_subscription" and the
// method that should be called to determine if the user is subscribing to a new channel
func DescriptorForYoutubeActionNewSubscription() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_subscription",
		Description: "When you subscribe to a new channel",
		Method:      isSubscribingToNewChannel,
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "string",
				Description: "The id of the channel you want to subscribe to",
				Required:    false,
			},
		},
		Components: []string{
			"youtube:sub:id",
			"youtube:sub:published_at",
			"youtube:sub:title",
			"youtube:sub:channel",
			"youtube:sub:description",
		},
	}
}
