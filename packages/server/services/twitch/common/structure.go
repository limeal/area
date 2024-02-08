package common

// `TwitchSoundTrack` is a struct that contains a `Track` struct that contains an array of `Artists`
// structs that contain a `ID`, `Name`, and `CreatorChannelID` string, a `ID`, `ISRC`, and `Duration`
// int, and a `Title` string, an `Album` struct that contains a `ID`, `Name`, and `ImageURL` string,
// and a `Source` struct that contains a `ID`, `ContentType`, `Title`, `ImageURL`, `SoundtrackURL`, and
// @property Track - The track that is currently playing.
// @property Source - The source of the music.
type TwitchSoundTrack struct {
	Track struct {
		Artists []struct {
			ID               string `json:"id"`
			Name             string `json:"name"`
			CreatorChannelID string `json:"creator_channel_id"`
		} `json:"artists"`
		ID       string `json:"id"`
		ISRC     string `json:"isrc"`
		Duration int    `json:"duration"`
		Title    string `json:"title"`
		Album    struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			ImageURL string `json:"image_url"`
		} `json:"album"`
	} `json:"track"`
	Source struct {
		ID            string `json:"id"`
		ContentType   string `json:"content_type"`
		Title         string `json:"title"`
		ImageURL      string `json:"image_url"`
		SoundtrackURL string `json:"soundtrack_url"`
		SpotifyURL    string `json:"spotify_url"`
	} `json:"source"`
}

// `TwitchSoundTracks` is a type that contains a slice of `TwitchSoundTrack`s.
// @property {[]TwitchSoundTrack} Data - An array of TwitchSoundTrack objects.
type TwitchSoundTracks struct {
	Data []TwitchSoundTrack `json:"data"`
}

// `TwitchStream` is a struct with a bunch of fields.
//
// The `json:"id"` part is called a struct tag. It tells the JSON encoder/decoder how to map the field
// to JSON.
//
// The `json:"id"` part is called a struct tag. It tells the JSON encoder/decoder how to map the field
// to JSON.
//
// The `json:"id"` part is called a struct tag. It tells the JSON encoder/decoder how to map the field
// to JSON.
//
// The `json:"id
// @property {string} ID - The ID of the stream.
// @property {string} UserID - The ID of the user who is streaming.
// @property {string} UserLogin - The user's login name.
// @property {string} UserName - The name of the user who is streaming
// @property {string} GameID - The ID of the game being played.
// @property {string} GameName - The name of the game being played
// @property {string} Type - The type of stream. Valid values: live, playlist, all. Default: live.
// @property {string} Title - The title of the stream.
// @property {[]string} Tags - A list of tags that have been applied to the stream.
// @property {int} ViewerCount - The number of viewers watching the stream.
// @property {string} StartedAt - The date and time when the stream started.
// @property {string} Language - The language of the stream.
// @property {string} ThumbnailURL - The URL of the thumbnail for the stream.
// @property {bool} IsMature - Whether the stream is flagged as mature or not.
type TwitchStream struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailURL string   `json:"thumbnail_url"`
	IsMature     bool     `json:"is_mature"`
}

// `TwitchStreams` is a type that contains a slice of `TwitchStream`s.
// @property {[]TwitchStream} Data - An array of TwitchStream objects.
type TwitchStreams struct {
	Data []TwitchStream `json:"data"`
}

// `TwitchClip` is a struct with a bunch of fields that are all strings except for `ViewCount` and
// `Duration` which are integers and floats respectively.
// @property {string} ID - The ID of the clip.
// @property {string} URL - The URL of the clip.
// @property {string} EmbedURL - The URL you can use to embed the clip in your own application.
// @property {string} BroadcasterID - The ID of the user who created the clip.
// @property {string} BroadcasterName - The name of the broadcaster who created the clip.
// @property {string} CreatorID - The ID of the user who created the clip.
// @property {string} CreatorName - The name of the user who created the clip.
// @property {string} VideoID - The ID of the video that the clip is taken from.
// @property {string} GameID - The ID of the game being played in the clip.
// @property {string} Language - The language of the clip.
// @property {string} Title - The title of the clip.
// @property {int} ViewCount - The number of times the clip has been viewed.
// @property {string} CreatedAt - The date and time when the clip was created.
// @property {string} ThumbnailURL - The URL of the thumbnail.
// @property {float64} Duration - The duration of the clip in seconds.
type TwitchClip struct {
	ID              string  `json:"id"`
	URL             string  `json:"url"`
	EmbedURL        string  `json:"embed_url"`
	BroadcasterID   string  `json:"broadcaster_id"`
	BroadcasterName string  `json:"broadcaster_name"`
	CreatorID       string  `json:"creator_id"`
	CreatorName     string  `json:"creator_name"`
	VideoID         string  `json:"video_id"`
	GameID          string  `json:"game_id"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	ViewCount       int     `json:"view_count"`
	CreatedAt       string  `json:"created_at"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	Duration        float64 `json:"duration"`
}

// `TwitchClips` is a type that contains a slice of `TwitchClip`s.
// @property {[]TwitchClip} Data - An array of TwitchClip objects.
type TwitchClips struct {
	Data []TwitchClip `json:"data"`
}

// `TwitchFollow` is a struct with 5 fields, all of which are strings.
// @property {string} FromId - The user ID of the user who followed the channel.
// @property {string} ToId - The ID of the channel that was followed.
// @property {string} ToLogin - The login name of the channel that was followed
// @property {string} ToName - The name of the channel that was followed
// @property {string} FollowedAt - The date and time when the user followed the channel.
type TwitchFollow struct {
	FromId     string `json:"from_id"`
	ToId       string `json:"to_id"`
	ToLogin    string `json:"to_login"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

// `TwitchFollows` is a type that contains a slice of `TwitchFollow`s.
// @property {[]TwitchFollow} Data - An array of TwitchFollow objects.
type TwitchFollows struct {
	Data []TwitchFollow `json:"data"`
}

// `TwitchGame` is a struct that contains a string called `ID`, a string called `Name`, a string called
// `BoxArtURL`, and a string called `IGDBID`.
// @property {string} ID - The ID of the game.
// @property {string} Name - The name of the game
// @property {string} BoxArtURL - The URL to the box art image for the game.
// @property {string} IGDBID - The ID of the game in the IGDB database.
type TwitchGame struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
	IGDBID    string `json:"igdb_id"`
}

// `TwitchGames` is a type that contains a slice of `TwitchGame`s.
// @property {[]TwitchGame} Data - An array of TwitchGame objects.
type TwitchGames struct {
	Data []TwitchGame `json:"data"`
}

// `TwitchUser` is a struct with a bunch of fields that are all strings except for `ViewCount` which is
// an int.
// @property {string} ID - The user's ID.
// @property {string} Login - The user's Twitch username.
// @property {string} DisplayName - The user's display name.
// @property {string} Type - The type of user. Valid values: staff, admin, global_mod, none
// @property {string} BroadcasterType - The broadcaster type, which will be one of "partner",
// "affiliate", or "".
// @property {string} Description - The user's channel description.
// @property {string} ProfileImageURL - The URL of the user's profile picture.
// @property {string} OfflineImageURL - The URL of the offline image of the stream.
// @property {int} ViewCount - The number of times the channel has been viewed.
// @property {string} Email - The user's email address. Returned if the request includes the
// user:read:email scope.
type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// `TwitchUsers` is a struct that contains a slice of `TwitchUser` structs.
// @property {[]TwitchUser} Data - An array of TwitchUser objects.
type TwitchUsers struct {
	Data []TwitchUser `json:"data"`
}
