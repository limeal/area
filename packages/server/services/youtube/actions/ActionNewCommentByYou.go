package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"
	"fmt"
)

// It gets the latest comment by the user on a video, channel or comment
func newCommentByYou(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)

	query["part"] = "snippet"
	query["maxResults"] = "1"
	query["textFormat"] = "plainText"

	switch (*req.Store)["req:entity:commentable:type"] {
	case "video":
		query["videoId"] = (*req.Store)["req:entity:commentable:id"].(string)
	case "channel":
		query["channelId"] = (*req.Store)["req:entity:commentable:id"].(string)
	case "comment":
		query["parentId"] = (*req.Store)["req:entity:commentable:id"].(string)
	}

	var err error
	var encode []byte

	switch (*req.Store)["req:entity:commentable:type"] {
	case "video", "channel":
		encode, _, err = req.Service.Endpoints["GetAllCommentsEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			query,
		})
	case "comment":
		encode, _, err = req.Service.Endpoints["GetAllRepliesEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			query,
		})
	}

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	comments := common.YoutubeCommentsResponse{}
	err = json.Unmarshal(encode, &comments)

	fmt.Println(comments.PageInfo.TotalResults)
	nbComments := comments.PageInfo.TotalResults
	if ok, errL := utils.IsLatestBasic(req.Store, nbComments); errL != nil || !ok {
		return shared.AreaResponse{Error: errL, Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"youtube:comment:id":            comments.Items[0].ID,
			"youtube:comment:like:count":    comments.Items[0].Snippet.LikeCount,
			"youtube:comment:viewer:rating": comments.Items[0].Snippet.ViewerRating,

			"youtube:video:id":   comments.Items[0].Snippet.VideoID,
			"youtube:channel:id": comments.Items[0].Snippet.ChannelID,

			"youtube:author:name":        comments.Items[0].Snippet.AuthorName,
			"youtube:author:channel:id":  comments.Items[0].Snippet.AuthorChannelID.Value,
			"youtube:author:channel:url": comments.Items[0].Snippet.AuthorChannel,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForYoutubeActionNewCommentByYou() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_comment_on_video",
		Description: "Triggered when a new comment is posted on a video",
		WIP:         true,
		RequestStore: map[string]static.StoreElement{
			"req:entity:commentable:type": {
				Description: "The entity type to comment",
				Type:        "select",
				Required:    true,
				Values: []string{
					"video",
					"channel",
					"comment",
				},
			},
			"req:entity:commentable:id": {
				Priority:    1,
				Description: "The entity id to comment (can be a videoId, channelId or commentId)",
				Type:        "string",
				Required:    true,
			},
		},
		Method: newCommentByYou,
		Components: []string{
			"youtube:comment:id",
			"youtube:comment:like:count",
			"youtube:comment:viewer:rating",
			"youtube:video:id",
			"youtube:channel:id",
			"youtube:author:name",
			"youtube:author:channel:id",
			"youtube:author:channel:url",
		},
	}
}
