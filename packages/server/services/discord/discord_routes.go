package discord

import (
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// It returns a slice of `static.ServiceRoute`s
func DiscordRoutes() []static.ServiceRoute {
	return []static.ServiceRoute{
		{
			Endpoint: "/guilds",
			Handler:  GetGuildsRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/guilds/:id/channels",
			Handler:  GetGuildChannelsRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/guilds/:id/members",
			Handler:  GetGuildMembersRoute,
			Method:   "GET",
			NeedAuth: true,
		},
		{
			Endpoint: "/channels/:id/messages",
			Handler:  GetChannelMessagesRoute,
			Method:   "GET",
			NeedAuth: true,
		},
	}
}

// ---------------------- Discord Routes ----------------------

// `UserGuildsResponse` is a struct with a field `Guilds` which is an array of structs with fields `ID`
// and `Name`.
// @property {[]struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 	}} Guilds - An array of guilds the user is in.
type UserGuildsResponse struct {
	Guilds []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"guilds"`
}

// It gets the guilds of the user or bot
func GetGuildsRoute(c *fiber.Ctx) error {

	authorization, errO := utils.VerifyRoute(c, "discord")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()
	bot := c.Query("bot", "")

	var encode []byte
	var errTw error

	if bot != "" {
		encode, _, errTw = service.Endpoints["GetBotGuildsEndpoint"].CallEncode([]interface{}{})
	} else {
		encode, _, errTw = service.Endpoints["GetUserGuildsEndpoint"].CallEncode([]interface{}{authorization})
	}

	if errTw != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errTw.Error(),
		})
	}

	userGuilds := UserGuildsResponse{}
	errrTh := json.Unmarshal(encode, &userGuilds.Guilds)
	if errrTh != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errrTh.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   userGuilds.Guilds,
			"fields": []string{"name", "id"},
		},
	})
}

// `GuildChannelsResponse` is a struct with a field `Channels` which is an array of structs with fields
// `ID` and `Name`.
// @property {[]struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 	}} Channels - An array of channels.
type GuildChannelsResponse struct {
	Channels []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channels"`
}

// It gets the guild channels of a guild
func GetGuildChannelsRoute(c *fiber.Ctx) error {

	auth0, errO := utils.VerifyRoute(c, "discord")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()
	id := c.Params("id")

	if id == "default" {
		authStore := make(map[string]interface{})
		err1 := json.Unmarshal(auth0.Other, &authStore)
		if err1 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err1.Error(),
			})
		}
		id = fmt.Sprintf("%v", authStore["guild_id"])
	}

	encode, _, errTw := service.Endpoints["ListGuildChannelsEndpoint"].CallEncode([]interface{}{id})
	if errTw != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errTw.Error(),
		})
	}

	channels := GuildChannelsResponse{}
	errrTh := json.Unmarshal(encode, &channels.Channels)
	if errrTh != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errrTh.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   channels.Channels,
			"fields": []string{"name", "id"},
		},
	})
}

// `Member` is a struct that contains a `User` struct that contains a `ID`, `Username`,
// `Discriminator`, and `Bot` field.
//
// Now that we have a type that matches the JSON, we can use the `json.Unmarshal` function to parse the
// JSON into a `Member` type.
// @property User - The user object of the member.
type Member struct {
	User struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
	} `json:"user"`
}

// `GuildMembersResponse` is a struct with a single field, `Members`, which is a slice of `Member`s.
// @property {[]Member} Members - An array of Member objects.
type GuildMembersResponse struct {
	Members []Member `json:"members"`
}

// It gets the guild members of a guild
func GetGuildMembersRoute(c *fiber.Ctx) error {

	auth0, errO := utils.VerifyRoute(c, "discord")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	service := Descriptor()
	id := c.Params("id")

	if id == "default" {
		authStore := make(map[string]interface{})
		err1 := json.Unmarshal(auth0.Other, &authStore)
		if err1 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  fiber.StatusInternalServerError,
				"error": err1.Error(),
			})
		}
		id = fmt.Sprintf("%v", authStore["guild_id"])
	}

	encode, _, errTw := service.Endpoints["ListGuildMembersEndpoint"].CallEncode([]interface{}{id})
	if errTw != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errTw.Error(),
		})
	}

	members := GuildMembersResponse{}
	errrTh := json.Unmarshal(encode, &members.Members)
	if errrTh != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errrTh.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   members.Members,
			"fields": []string{"user:username", "user:id"},
		},
	})
}

// `Message` is a struct that contains a string `ID`, a string `ChannelID`, an `Author` struct, and a
// string `Content`.
//
// The `Author` struct contains a string `ID`, a string `Username`, a string `Discriminator`, and a
// boolean `Bot`.
// @property {string} ID - The ID of the message.
// @property {string} ChannelID - The ID of the channel the message was sent in.
// @property Author - The author of the message.
// @property {string} Content - The message content.
type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Author    struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
	} `json:"author"`
	Content string `json:"content"`
}

// `ChannelMessagesResponse` is a struct with a field `Messages` of type `[]Message`.
// @property {[]Message} Messages - An array of messages.
type ChannelMessagesResponse struct {
	Messages []Message `json:"messages"`
}

// It gets the messages from a channel
func GetChannelMessagesRoute(c *fiber.Ctx) error {
	limitP := c.Params("limit", "5")

	_, errO := utils.VerifyRoute(c, "discord")
	if errO != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":  fiber.StatusForbidden,
			"error": "Forbidden",
		})
	}

	limit, err := strconv.Atoi(limitP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	service := Descriptor()
	encode, _, errTw := service.Endpoints["ListChannelMessagesEndpoint"].CallEncode([]interface{}{
		c.Params("id"),
		map[string]string{
			"limit": "20",
		},
	})
	if errTw != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errTw.Error(),
		})
	}

	channelMessages := ChannelMessagesResponse{}
	errrTh := json.Unmarshal(encode, &channelMessages.Messages)
	if errrTh != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":  fiber.StatusInternalServerError,
			"error": errrTh.Error(),
		})
	}

	var messages []Message
	if len(channelMessages.Messages) > limit {
		messages = channelMessages.Messages[:limit]
	} else {
		messages = channelMessages.Messages
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"data":   messages,
			"fields": []string{"content", "id"},
		},
	})
}
