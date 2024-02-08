package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/reddit/common"
	"area-server/utils"
	"encoding/json"
)

// It checks if the number of comments by the user is the same as the last time the function was called
func newCommentByYou(req static.AreaRequest) shared.AreaResponse {

	user := req.AuthStore["name"].(string)

	query := make(map[string]string)
	query["limit"] = "100"
	query["sort"] = "new"

	encode, _, err := req.Service.Endpoints["GetUserCommentsEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		user,
		query,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	newComment := common.NewCommentByYouResponse{}
	errr := json.Unmarshal(encode, &newComment)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbComment := len(newComment.Data.Children)
	if ok, errL := utils.IsLatestBasic(req.Store, nbComment); errL != nil || !ok {
		return shared.AreaResponse{Error: errL, Success: false}
	}

	return shared.AreaResponse{Success: true}
}

// `newCommentByYou` is a function that takes a `reddit.NewCommentByYou` struct and returns a
// `static.ServiceArea` struct
func DescriptorForRedditActionNewCommentByYou() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_comment_by_you",
		Description: "Triggered when you post a new comment",
		Method:      newCommentByYou,
	}
}
