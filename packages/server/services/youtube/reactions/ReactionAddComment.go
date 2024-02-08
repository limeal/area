package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// A ThreadCommentSnippet is a struct that contains a VideoId, a ChannelId, and a TopLevelComment,
// which is a struct that contains a Snippet, which is a struct that contains a TextOriginal.
// @property {string} VideoId - The id of the video that the comment is on.
// @property {string} ChannelId - The ID of the channel that the comment was posted to.
// @property TopLevelComment - The top-level comment in the thread.
type ThreadCommentSnippet struct {
	VideoId         string `json:"videoId"`
	ChannelId       string `json:"channelId"`
	TopLevelComment struct {
		Snippet struct {
			TextOriginal string `json:"textOriginal"`
		} `json:"snippet"`
	} `json:"topLevelComment"`
}

// `CommentSnippet` is a struct with two fields, `ParentId` and `TextOriginal`.
//
// The `json` tags are used to tell the `encoding/json` package how to map the fields to JSON.
//
// The `ParentId` field is a string, and the `TextOriginal` field is a string.
//
// The `ParentId` field is mapped to the `commentId` key in the JSON.
//
// The `TextOriginal` field is mapped to the `textOriginal` key in the JSON.
// @property {string} ParentId - The ID of the parent comment.
// @property {string} TextOriginal - The comment's text.
type CommentSnippet struct {
	ParentId     string `json:"commentId"`
	TextOriginal string `json:"textOriginal"`
}

// It adds a comment to a video, channel or comment
func addComment(req static.AreaRequest) shared.AreaResponse {

	body := make(map[string]interface{})

	commentText := utils.GenerateFinalComponent((*req.Store)["req:comment:text"].(string), req.ExternalData, []string{})
	commentTableID := utils.GenerateFinalComponent((*req.Store)["req:entity:commentable:id"].(string), req.ExternalData, []string{
		"youtube:comment:id",
		"youtube:video:id",
		"youtube:channel:id",
	})

	switch (*req.Store)["req:entity:commentable:type"] {
	case "video":
		body["snippet"] = ThreadCommentSnippet{
			VideoId: commentTableID,
			TopLevelComment: struct {
				Snippet struct {
					TextOriginal string "json:\"textOriginal\""
				} "json:\"snippet\""
			}{
				Snippet: struct {
					TextOriginal string "json:\"textOriginal\""
				}{
					TextOriginal: commentText,
				},
			},
		}
	case "channel":
		body["snippet"] = ThreadCommentSnippet{
			ChannelId: commentTableID,
			TopLevelComment: struct {
				Snippet struct {
					TextOriginal string "json:\"textOriginal\""
				} "json:\"snippet\""
			}{
				Snippet: struct {
					TextOriginal string "json:\"textOriginal\""
				}{
					TextOriginal: commentText,
				},
			},
		}
	case "comment":
		body["snippet"] = CommentSnippet{
			ParentId:     commentTableID,
			TextOriginal: commentText,
		}
	}

	ebody, err := json.Marshal(body)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var errr error
	switch (*req.Store)["req:entity:commentable:type"] {
	case "video", "channel":
		_, _, errr = req.Service.Endpoints["AddCommentThreadEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			string(ebody),
		})
	case "comment":
		_, _, errr = req.Service.Endpoints["AddCommentEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			string(ebody),
		})
	}

	return shared.AreaResponse{
		Error: errr,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForYoutubeReactionAddComment() static.ServiceArea {
	return static.ServiceArea{
		Name:        "add_comment",
		Description: "Add a comment to a video",
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
			"req:comment:text": {
				Priority:    2,
				Description: "The text of the comment",
				Type:        "string",
				Required:    true,
			},
		},
		Method: addComment,
	}
}
