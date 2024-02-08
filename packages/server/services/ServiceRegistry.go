package services

import (
	"area-server/classes/static"
	"area-server/services/discord"
	"area-server/services/dropbox"
	"area-server/services/github"
	"area-server/services/gmail"
	"area-server/services/openweather"
	"area-server/services/reddit"
	"area-server/services/spotify"
	"area-server/services/time"
	"area-server/services/twitch"
	"area-server/services/webhook"
	"area-server/services/youtube"
)

// It's a list of all the services that are available.
var List []static.Service = []static.Service{
	discord.Descriptor(),     // Discord
	dropbox.Descriptor(),     // Dropbox
	github.Descriptor(),      // Github
	gmail.Descriptor(),       // Gmail
	openweather.Descriptor(), // Open Weather
	reddit.Descriptor(),      // Reddit
	spotify.Descriptor(),     // Spotify
	time.Descriptor(),        // Time
	twitch.Descriptor(),      // Twitch
	webhook.Descriptor(),     // Webhook
	youtube.Descriptor(),     // YouTube
}

// GetServiceByName returns a pointer to a static.Service struct if the name of the service matches the
// name passed in as an argument.
func GetServiceByName(name string) *static.Service {
	for _, service := range List {
		if service.Name == name {
			return &service
		}
	}
	return nil
}
