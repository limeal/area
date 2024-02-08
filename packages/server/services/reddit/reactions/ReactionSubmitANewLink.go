package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
)

// It takes a request, and returns a response
func submitNewLink(req static.AreaRequest) shared.AreaResponse {

	subReddit := utils.GenerateFinalComponent((*req.Store)["req:subreddit:name"].(string), req.ExternalData, []string{
		"reddit:subreddit:name",
	})
	title := utils.GenerateFinalComponent((*req.Store)["req:link:title"].(string), req.ExternalData, []string{})
	url := utils.GenerateFinalComponent((*req.Store)["req:link:url"].(string), req.ExternalData, []string{
		".+url",
	})

	_, _, err := req.Service.Endpoints["SubmitNewLinkEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		subReddit,
		title,
		url,
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a `static.ServiceArea` object that describes the service area
func DescriptorForRedditReactionSubmitNewLink() static.ServiceArea {
	return static.ServiceArea{
		Name:        "submit_new_link",
		Description: "Submit a new link to Reddit",
		RequestStore: map[string]static.StoreElement{
			"req:subreddit:name": {
				Description: "The subreddit to submit the link to",
				Required:    true,
			},
			"req:link:title": {
				Priority:    1,
				Description: "The title of the link",
				Required:    true,
			},
			"req:link:url": {
				Priority:    2,
				Description: "The URL of the link",
				Required:    true,
			},
		},
		Method: submitNewLink,
	}
}
