package discord

import (
	"area-server/classes/static"
	"area-server/db/postgres/models"
	"area-server/services/common"
	"regexp"

	"github.com/forPelevin/gomoji"
)

// It returns a map of validators for the Discord service
func DiscordValidators() static.ServiceValidator {
	return static.ServiceValidator{
		"req:guild:id":           GuildIDValidator,
		"req:channel:id":         ChannelIDValidator,
		"req:user:id":            UserIDValidator,
		"req:allow:bot":          common.BoolValidator,
		"req:not:reply":          common.BoolValidator,
		"req:message:regex":      RegexValidator,
		"req:emoji":              EmojiValidator,
		"req:gateway:event:type": GatewayEventTypeValidator,
	}
}

// ------------------------- Validators ------------------------------

// It checks if the value is a string, and if it is, it checks if the value is a valid guild ID
func GuildIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value.(string) == "{{discord:guild:id}}" {
		return true
	}

	_, _, err := service.Endpoints["GetGuildEnpoint"].CallEncode([]interface{}{value.(string)})
	if err != nil {
		return false
	}

	return true
}

// It checks if the value is a string, and if it is, it checks if the value is the special string
// `{{discord:channel:id}}`, and if it isn't, it checks if the value is a valid channel ID
func ChannelIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value.(string) == "{{discord:channel:id}}" {
		return true
	}

	_, _, err := service.Endpoints["GetChannelEndpoint"].CallEncode([]interface{}{value.(string)})
	if err != nil {
		return false
	}

	return true
}

// If the value is a string, and the value is equal to the string `{{discord:user:id}}`, or if the
// value is a string and the Discord API returns a user with the ID of the value, then the value is
// valid
func UserIDValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value.(string) == "{{discord:user:id}}" {
		return true
	}

	_, _, err := service.Endpoints["GetUserEndpoint"].CallEncode([]interface{}{value.(string)})
	if err != nil {
		return false
	}

	return true
}

// It checks if the value is a string and if it is, it checks if it's a valid regular expression
func RegexValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if _, err := regexp.Compile(value.(string)); err != nil {
		return false
	}

	return true
}

// `EmojiValidator` checks if the value is a string and if it is, it checks if it's a valid emoji
func EmojiValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	if value.(string) == "{{discord:emoji}}" {
		return true
	}

	result := gomoji.FindAll(value.(string))
	if result != nil && len(result) == 1 {
		return true
	}

	return false
}

// It checks if the value is a string, and if it is, it checks if it's one of the valid gateway event
// types
func GatewayEventTypeValidator(authorization *models.Authorization, service *static.Service, value interface{}, store map[string]interface{}) bool {
	if value == nil {
		return false
	}

	if _, ok := value.(string); !ok {
		return false
	}

	switch value.(string) {
	case "READY", "RESUMED", "CHANNEL_CREATE", "CHANNEL_UPDATE", "CHANNEL_DELETE", "CHANNEL_PINS_UPDATE", "GUILD_CREATE", "GUILD_UPDATE", "GUILD_DELETE", "GUILD_BAN_ADD", "GUILD_BAN_REMOVE", "GUILD_EMOJIS_UPDATE", "GUILD_INTEGRATIONS_UPDATE", "GUILD_MEMBER_ADD", "GUILD_MEMBER_REMOVE", "GUILD_MEMBER_UPDATE", "GUILD_MEMBERS_CHUNK", "GUILD_ROLE_CREATE", "GUILD_ROLE_UPDATE", "GUILD_ROLE_DELETE", "MESSAGE_CREATE", "MESSAGE_UPDATE", "MESSAGE_DELETE", "MESSAGE_DELETE_BULK", "MESSAGE_REACTION_ADD", "MESSAGE_REACTION_REMOVE", "MESSAGE_REACTION_REMOVE_ALL", "PRESENCE_UPDATE", "TYPING_START", "USER_UPDATE", "VOICE_STATE_UPDATE", "VOICE_SERVER_UPDATE", "WEBHOOKS_UPDATE":
		return true
	}
	return false
}
