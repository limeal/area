package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/services/youtube/common"
	"area-server/utils"
	"encoding/json"
)

// It checks if the user has liked a new video
func hasLikedANewVideo(req static.AreaRequest) shared.AreaResponse {
	query := make(map[string]string)

	query["part"] = "snippet"
	query["maxResults"] = "1"
	query["myRating"] = "like"

	encode, _, err := req.Service.Endpoints["GetAllVideosEndpoint"].CallEncode([]interface{}{req.Authorization, query})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	videoList := common.YoutubeVideoResponse{}
	errr := json.Unmarshal(encode, &videoList)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbLikedVideos := videoList.PageInfo.TotalResults
	ok, errL := utils.IsLatestBasic(req.Store, nbLikedVideos)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"youtube:video:id":           videoList.Items[0].ID,
			"youtube:video:title":        videoList.Items[0].Snippet.Title,
			"youtube:video:description":  videoList.Items[0].Snippet.Description,
			"youtube:video:channel":      videoList.Items[0].Snippet.ChannelTitle,
			"youtube:video:published_at": videoList.Items[0].Snippet.PublishedAt,
		},
	}
}

// `DescriptorForYoutubeActionNewLikedVideo` returns a `static.ServiceArea` that describes the
// `new_liked_video` action for the `youtube` service
func DescriptorForYoutubeActionNewLikedVideo() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_liked_video",
		Description: "When a new video is liked",
		Method:      hasLikedANewVideo,
		Components: []string{
			"youtube:video:id",
			"youtube:video:title",
			"youtube:video:description",
			"youtube:video:channel",
			"youtube:video:published_at",
		},
	}
}
