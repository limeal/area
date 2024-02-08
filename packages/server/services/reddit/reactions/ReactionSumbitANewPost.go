package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
)

// It takes a request, generates a title and content from the request's store and external data, and
// then calls the Reddit API to submit a new post
func submitNewPost(req static.AreaRequest) shared.AreaResponse {

	subReddit := utils.GenerateFinalComponent((*req.Store)["req:subreddit:name"].(string), req.ExternalData, []string{
		"reddit:subreddit:name",
	})
	title := utils.GenerateFinalComponent((*req.Store)["req:post:title"].(string), req.ExternalData, []string{})
	content := utils.GenerateFinalComponent((*req.Store)["req:post:body"].(string), req.ExternalData, []string{})

	_, _, err := req.Service.Endpoints["SubmitNewPostEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		subReddit,
		title,
		content,
	})

	return shared.AreaResponse{
		Error: err,
	}
}

// It returns a `static.ServiceArea` object that describes the service area
func DescriptorForRedditReactionSubmitNewPost() static.ServiceArea {
	return static.ServiceArea{
		Name:        "submit_new_post",
		Description: "Submit a new post to Reddit",
		RequestStore: map[string]static.StoreElement{
			"req:subreddit:name": {
				Description: "The subreddit to submit the link to",
				Required:    true,
			},
			"req:post:title": {
				Priority:    1,
				Description: "The title of the post",
				Required:    true,
			},
			"req:post:body": {
				Priority:    2,
				Type:        "long_string",
				Description: "The body of the post",
				Required:    true,
			},
		},
		Method: submitNewPost,
	}
}
