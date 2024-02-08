package discord

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/utils"
	"os"
)

// It returns a map of all the endpoints that the Discord API has
func DiscordEndpoints() static.ServiceEndpoint {
	return static.ServiceEndpoint{
		// Routes
		"GetBotGuildsEndpoint": {
			BaseURL:        "https://discord.com/api/v10/users/@me/guilds",
			Params:         GetBotGuildsEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Validations
		"GetGuildEnpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetChannelEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetUserEndpoint": {
			BaseURL:        "https://discord.com/api/v10/users/${id}",
			Params:         GetBasicEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetChannelMessageEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${cid}/messages/${id}",
			Params:         GetChannelMessageEndpointParams,
			ExpectedStatus: []int{200},
		},
		"GetGuildEmojiEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${gid}/emojis/${id}",
			Params:         GetGuildEmojiEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Actions
		"ListChannelMessagesEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}/messages",
			Params:         GetChannelsEndpointParams,
			ExpectedStatus: []int{200, 429},
		},
		"GetUserGuildsEndpoint": {
			BaseURL:        "https://discord.com/api/v10/users/@me/guilds",
			Params:         GetUserGuildsEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListGuildMembersEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}/members",
			Params:         ListGuildMembersEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListGuildChannelsEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}/channels",
			Params:         ListGuildEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListGuildEventsEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}/scheduled-events",
			Params:         ListGuildEndpointParams,
			ExpectedStatus: []int{200},
		},
		"ListPinnedMessagesEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}/pins",
			Params:         ListGuildEndpointParams,
			ExpectedStatus: []int{200, 429},
		},
		"ListUserThatReactedEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}/messages/${mid}/reactions/${emoji}",
			Params:         ListUserThatReactedEndpointParams,
			ExpectedStatus: []int{200},
		},
		// Reactions
		"CreateDMEndpoint": {
			BaseURL:        "https://discord.com/api/v10/users/@me/channels",
			Params:         CreateDMEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateMessageEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}/messages",
			Params:         CreateMessageEndpointParams,
			ExpectedStatus: []int{200},
		},
		"CreateReactionEndpoint": {
			BaseURL:        "https://discord.com/api/v10/channels/${id}/messages/${message_id}/reactions/${emoji}/@me",
			Params:         CreateReactionEndpointParams,
			ExpectedStatus: []int{204},
		},
		"BanUserEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}/bans/${uid}",
			Params:         BanUserEndpointParams,
			ExpectedStatus: []int{204},
		},
		"KickUserEndpoint": {
			BaseURL:        "https://discord.com/api/v10/guilds/${id}/members/${uid}",
			Params:         KickUserEndpointParams,
			ExpectedStatus: []int{204},
		},
	}
}

// ------------------------- Endpoint Params ------------------------------

// It returns a pointer to a `RequestParams` struct with the method set to `GET` and the
// `Authorization` header set to `Bot <DISCORD_BOT_TOKEN>`
func GetBotGuildsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
	}
}

// It takes a slice of interfaces, and returns a pointer to a RequestParams struct
func GetBasicEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
	}
}

// It returns a pointer to a RequestParams struct with the method, headers, and url params set
func GetChannelMessageEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"cid": params[0].(string),
			"id":  params[1].(string),
		},
	}
}

// `GetGuildEmojiEndpointParams` returns a `RequestParams` struct with the method `GET`, the headers
// `Authorization` and `Bot <DISCORD_BOT_TOKEN>`, and the url params `gid` and `id`
func GetGuildEmojiEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"gid": params[0].(string),
			"id":  params[1].(string),
		},
	}
}

// It creates a request params object for the Discord API endpoint that sends a message to a channel
func CreateMessageEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
		Body: params[1].(string),
	}
}

// It creates a request params object for the endpoint that adds a reaction to a message
func CreateReactionEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id":         params[0].(string),
			"message_id": params[1].(string),
			"emoji":      params[2].(string),
		},
	}
}

// It takes in a slice of interfaces, and returns a pointer to a RequestParams struct
func GetChannelsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
		QueryParams: params[1].(map[string]string),
	}
}

// It takes a slice of interfaces, checks if the first element is a pointer to a `models.Authorization`
// struct, and if it is, it returns a `utils.RequestParams` struct with the `Authorization` header set
// to the `AccessToken` field of the `models.Authorization` struct
func GetUserGuildsEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + params[0].(*models.Authorization).AccessToken,
		},
	}
}

// It takes a slice of interfaces, and returns a pointer to a RequestParams struct
func ListGuildMembersEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
		QueryParams: map[string]string{
			"limit": "1000",
		},
	}
}

// It returns a pointer to a RequestParams struct with the method, headers, and url params set
func ListGuildEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
	}
}

// It takes in two parameters, the first being the guild ID and the second being the search query. It
// then returns a request object with the method set to GET, the headers set to the authorization and
// content type, the URL parameters set to the guild ID, and the body set to the search query
func SearchGuildMembersEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
			"Content-Type":  "application/json",
		},
		UrlParams: map[string]string{
			"id": params[0].(string),
		},
		Body: "{\"query\": \"" + params[1].(string) + "\"}",
	}
}

// `ListUserThatReactedEndpointParams` takes in a slice of `interface{}` and returns a
// `*utils.RequestParams`
//
// The first thing you'll notice is that the function takes in a slice of `interface{}`. This is
// because the function is generic and can be used for any endpoint that takes in a slice of
// `interface{}` as parameters
func ListUserThatReactedEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id":    params[0].(string),
			"mid":   params[1].(string),
			"emoji": params[2].(string),
		},
		QueryParams: params[3].(map[string]string),
	}
}

// It returns a pointer to a RequestParams struct with the method, headers, and url params set
func BanUserEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "PUT",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id":  params[0].(string),
			"uid": params[1].(string),
		},
	}
}

// It returns a pointer to a RequestParams struct with the method, headers, and url parameters set
func KickUserEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "DELETE",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
		},
		UrlParams: map[string]string{
			"id":  params[0].(string),
			"uid": params[1].(string),
		},
	}
}

// It creates a request params object for the Discord API endpoint that creates a direct message
// channel between the bot and a user
func CreateDMEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
			"Content-Type":  "application/json",
		},
		Body: "{\"recipient_id\": \"" + params[0].(string) + "\"}",
	}
}
