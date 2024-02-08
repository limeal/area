package common

// `SubReddit` is a type that represents a subreddit.
// @property {string} DisplayName - The name of the subreddit
// @property {string} ID - The unique ID of the subreddit.
// @property {string} Title - The title of the subreddit
// @property {string} URL - The URL of the subreddit.
// @property {int} Subscribers - The number of subscribers to the subreddit.
// @property {string} IconImg - The URL of the subreddit's icon.
// @property {string} Description - A short blurb describing the subreddit.
type SubReddit struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Subscribers int    `json:"subscribers"`
	IconImg     string `json:"icon_img"`
	Description string `json:"public_description"`
}

// `SubRedditResponse` is a struct with a `Kind` field of type `string` and a `Data` field of type
// `struct` with `Children` field of type `[]struct` with `Kind` and `Data` fields of type `string` and
// `SubReddit` respectively.
// @property {string} Kind - The type of data that is being returned.
// @property Data - This is the main data structure that contains the subreddit data.
type SubRedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string    `json:"kind"`
			Data SubReddit `json:"data"`
		} `json:"children"`
		After string `json:"after"`
	} `json:"data"`
}

// User is a type that has a field called Name that is a string and is exported.
// @property {string} Name - The name of the user.
type User struct {
	Name string `json:"name"`
}

// A RedditPost is a struct that contains the fields ID, Name, Title, Author, URL, Text, SubRedditID,
// SubRedditName, and SubRedditSubscribers.
// @property {string} ID - The unique ID of the post.
// @property {string} Name - The unique ID of the post.
// @property {string} Title - The title of the post
// @property {string} Author - The username of the author of the post.
// @property {string} URL - The URL of the post
// @property {string} Text - The text of the post.
// @property {string} SubRedditID - The ID of the subreddit the post is in.
// @property {string} SubRedditName - The name of the subreddit the post is in
// @property {int} SubRedditSubscribers - The number of subscribers to the subreddit
type RedditPost struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Title                string `json:"title"`
	Author               string `json:"author"`
	URL                  string `json:"url"`
	Text                 string `json:"selftext"`
	SubRedditID          string `json:"subreddit_id"`
	SubRedditName        string `json:"subreddit_name_prefixed"`
	SubRedditSubscribers int    `json:"subreddit_subscribers"`
}

// A NewPostsOnSubRedditResponse is a struct that contains a string called Kind and a struct called
// Data that contains a slice of structs called Children that contain a string called Kind and a
// RedditPost called Data.
// @property {string} Kind - The type of data that is being returned.
// @property Data - This is the main data object that contains the posts.
type NewPostsOnSubRedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string     `json:"kind"`
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// `RedditComment` is a struct that contains a string for each of the following fields: `ID`,
// `LinkTitle`, `Author`, `Body`, `SubReddit`, and `SubRedditID`.
// @property {string} ID - The ID of the comment
// @property {string} LinkTitle - The title of the post that the comment is on
// @property {string} Author - The author of the comment
// @property {string} Body - The text of the comment
// @property {string} SubReddit - The name of the subreddit the comment was posted in
// @property {string} SubRedditID - The ID of the subreddit the comment was posted in.
type RedditComment struct {
	ID          string `json:"id"`
	LinkTitle   string `json:"link_title"`
	Author      string `json:"link_author"`
	Body        string `json:"body"`
	SubReddit   string `json:"subreddit"`
	SubRedditID string `json:"subreddit_id"`
}

// NewCommentByYouResponse is a struct that contains a string and a struct that contains a slice of
// structs that contain a string and a RedditComment struct.
// @property {string} Kind - The type of data that is being returned.
// @property Data - This is the main data object that contains the actual comment.
type NewCommentByYouResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string        `json:"kind"`
			Data RedditComment `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
