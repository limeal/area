package common

// `YoutubeVideoResponse` is a struct that contains a slice of structs that contain a string, a struct
// that contains a string, a string, a string, a string, and a string, a struct that contains an int
// and an int.
// @property {[]struct {
// 		ID      string `json:"id"`
// 		Snippet struct {
// 			PublishedAt  string `json:"publishedAt"`
// 			ChannelID    string `json:"channelId"`
// 			Title        string `json:"title"`
// 			Description  string `json:"description"`
// 			ChannelTitle string `json:"channelTitle"`
// 		} `json:"snippet"`
// 	}} Items - This is an array of videos that match the search criteria.
// @property PageInfo - This is the information about the page of results that you're currently looking
// at.
type YoutubeVideoResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelID    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// It's a struct that contains a slice of structs that contain a struct that contains a struct that
// contains a struct that contains a struct that contains a struct that contains a struct that contains
// a struct that contains a struct that contains a struct that contains a struct that contains a struct
// that contains a struct that contains a struct that contains a struct that contains a struct that
// contains a struct that contains a struct that contains a struct that contains a struct that contains
// a struct that contains a struct that contains a struct that contains a struct that contains a struct
// that contains a struct that contains a struct that contains a struct that
// @property {[]struct {
// 		ID      string `json:"id"`
// 		Snippet struct {
// 			AuthorName      string `json:"authorDisplayName"`
// 			AuthorChannel   string `json:"authorChannelUrl"`
// 			AuthorChannelID struct {
// 				Value string `json:"value"`
// 			} `json:"authorChannelId"`
// 			ChannelID    string `json:"channelId"`
// 			VideoID      string `json:"videoId"`
// 			LikeCount    int    `json:"likeCount"`
// 			ViewerRating int    `json:"viewerRating"`
// 		}
// 	}} Items - This is an array of comments.
// @property PageInfo - This is the information about the page of results that you're currently
// viewing.
type YoutubeCommentsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			AuthorName      string `json:"authorDisplayName"`
			AuthorChannel   string `json:"authorChannelUrl"`
			AuthorChannelID struct {
				Value string `json:"value"`
			} `json:"authorChannelId"`
			ChannelID    string `json:"channelId"`
			VideoID      string `json:"videoId"`
			LikeCount    int    `json:"likeCount"`
			ViewerRating int    `json:"viewerRating"`
		}
	}
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

type YoutubeSubscriptionResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelID   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// `YoutubePlaylistsResponse` is a struct that contains a slice of structs that contain a string and a
// struct that contains a string and a string and a string and a string and a string.
// @property {[]struct {
// 		ID      string `json:"id"`
// 		Snippet struct {
// 			PublishedAt  string `json:"publishedAt"`
// 			ChannelID    string `json:"channelId"`
// 			Title        string `json:"title"`
// 			Description  string `json:"description"`
// 			ChannelTitle string `json:"channelTitle"`
// 		} `json:"snippet"`
// 	}} Items - An array of playlist resources that match the request criteria.
// @property PageInfo - This is a struct that contains the total number of results and the number of
// results per page.
type YoutubePlaylistsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelID    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// `YoutubePlaylistItemsResponse` is a struct that contains a slice of structs that contain a string, a
// string, and a struct that contains a string and a string.
// @property {[]struct {
// 		ID      string `json:"id"`
// 		Snippet struct {
// 			PublishedAt  string `json:"publishedAt"`
// 			ChannelID    string `json:"channelId"`
// 			Title        string `json:"title"`
// 			Description  string `json:"description"`
// 			ChannelTitle string `json:"channelTitle"`
// 			PlaylistID   string `json:"playlistId"`
// 			ResourceID   struct {
// 				Kind    string `json:"kind"`
// 				VideoID string `json:"videoId"`
// 			} `json:"resourceId"`
// 		} `json:"snippet"`
// 	}} Items - An array of playlist items that match the request criteria.
// @property PageInfo - This is a struct that contains the total number of results and the number of
// results per page.
type YoutubePlaylistItemsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelID    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			PlaylistID   string `json:"playlistId"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// "YoutubeChannelsResponse is a struct that contains a slice of structs that contain a struct that
// contains a struct that contains a struct that contains a struct that contains a struct that contains
// a struct that contains a struct that contains a struct that contains a struct that contains a struct
// that contains a struct that contains a struct that contains a struct that contains a struct that
// contains a struct that contains a struct that contains a struct that contains a struct that contains
// a struct that contains a struct that contains a struct that contains a struct that contains a struct
// that contains a struct that contains a struct that contains a struct
// @property {[]struct {
// 		ID      string `json:"id"`
// 		Snippet struct {
// 			PublishedAt  string `json:"publishedAt"`
// 			ChannelID    string `json:"channelId"`
// 			Title        string `json:"title"`
// 			Description  string `json:"description"`
// 			ChannelTitle string `json:"channelTitle"`
// 		} `json:"snippet"`
// 	}} Items - An array of channel objects that match the request criteria. Each item object contains
// information about the channel, such as its title, description, and thumbnail images.
// @property PageInfo - This is a struct that contains the total number of results and the number of
// results per page.
type YoutubeChannelsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt  string `json:"publishedAt"`
			ChannelID    string `json:"channelId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}
