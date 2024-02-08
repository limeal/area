package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"
)

// It checks if there's a new subscriber to the channel
func hasNewSubscriber(req static.AreaRequest) shared.AreaResponse {
	query := make(map[string]string)

	query["part"] = "subscriberSnippet"
	query["maxResults"] = "1"
	query["myRecentSubscribers"] = "true"

	encode, _, err := req.Service.Endpoints["GetAllSubscriptionsEndpoint"].CallEncode([]interface{}{req.Authorization, query})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	subList := common.YoutubeSubscriptionResponse{}
	errr := json.Unmarshal(encode, &subList)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbSubscribers := subList.PageInfo.TotalResults
	ok, errL := utils.IsLatestBasic(req.Store, nbSubscribers)
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

// When you have a new subscriber
func DescriptorForYoutubeActionNewSubscriber() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_subscriber",
		Description: "When you have a new subscriber",
		RequestStore: map[string]static.StoreElement{
			"req:channel:id": {
				Type:        "string",
				Description: "The channel id of the new subscriber",
				Required:    false,
			},
		},
		Method: hasNewSubscriber,
		Components: []string{
			"youtube:sub:id",
			"youtube:sub:published_at",
			"youtube:sub:title",
			"youtube:sub:channel",
			"youtube:sub:description",
		},
	}
}
