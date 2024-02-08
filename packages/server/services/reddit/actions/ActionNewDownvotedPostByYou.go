package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/reddit/common"
	"area-server/utils"
	"encoding/json"
)

// It gets the latest downvoted post by the user, and returns it
func newDownvotedPostByYou(req static.AreaRequest) shared.AreaResponse {

	userName := req.AuthStore["name"].(string)

	query := make(map[string]string)
	query["limit"] = "100"
	query["sort"] = "new"

	encode, _, err := req.Service.Endpoints["GetUserDownvotedPostsEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		userName,
		query,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	newPosts := common.NewPostsOnSubRedditResponse{}
	errr := json.Unmarshal(encode, &newPosts)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbPosts := len(newPosts.Data.Children)
	ok, errL := utils.IsLatestBasic(req.Store, nbPosts)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			// Subreddit
			"reddit:subreddit:id":   newPosts.Data.Children[0].Data.SubRedditID,
			"reddit:subreddit:name": newPosts.Data.Children[0].Data.SubRedditName,
			// Post
			"reddit:post:id":     newPosts.Data.Children[0].Data.ID,
			"reddit:post:name":   newPosts.Data.Children[0].Data.Name,
			"reddit:post:title":  newPosts.Data.Children[0].Data.Title,
			"reddit:post:author": newPosts.Data.Children[0].Data.Author,
			"reddit:post:url":    newPosts.Data.Children[0].Data.URL,
			"reddit:post:text":   newPosts.Data.Children[0].Data.Text,
		},
	}
}

// `newDownvotedPostByYou` is a function that takes a `reddit.Subreddit` and a `reddit.Post` and
// returns a `static.ServiceArea`
//
// The `static.ServiceArea` is a struct that contains the following fields:
//
// - `Name`: The name of the service area. This is used to identify the service area in the database.
// - `Description`: A description of the service area. This is used to describe the service area in the
// UI.
// - `Method`: The function that will be called when the service area is triggered.
// - `Components`: A list of components that are required to be present in the service area
func DescriptorForRedditActionNewDownvotedPostByYou() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_downvoted_post_by_you",
		Description: "Triggered when you downvote a post (<100)",
		Method:      newDownvotedPostByYou,
		Components: []string{
			// Subreddit
			"reddit:subreddit:id",
			"reddit:subreddit:name",
			// Post
			"reddit:post:id",
			"reddit:post:name",
			"reddit:post:title",
			"reddit:post:author",
			"reddit:post:url",
			"reddit:post:text",
		},
	}
}
