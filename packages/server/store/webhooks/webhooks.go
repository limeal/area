package webhooks

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
)

type WebhookMode int

// It's defining a new type called WebhookMode, and then defining two constants of that type.
const (
	WebhookTypeAppletTrigger WebhookMode = iota
	WebhookTypeServiceInteraction
)

// `Webhook` is a struct that has two fields, `Mode` and `Data`. `Mode` is a `WebhookMode` and `Data`
// is a channel of `[]byte`.
// @property {WebhookMode} Mode - The mode of the webhook. This can be either "subscribe" or
// "unsubscribe".
// @property Data - This is the channel that the webhook will send data to.
type Webhook struct {
	Mode WebhookMode
	Data (chan []byte)
}

// `HistoryItem` is a struct with two fields, `ID` and `Author`, both of which are strings.
// @property {string} ID - The ID of the item.
// @property {string} Author - The name of the person who made the change.
type HistoryItem struct {
	ID     string `json:"id"`
	Author string `json:"author"`
}

var History = list.New()
var Webhooks = make(map[string]Webhook)

// It adds a webhook to the Webhooks map
func AddWebhook(webhookName string) {
	fmt.Println("Adding webhook: " + webhookName)
	Webhooks[webhookName] = Webhook{
		Mode: WebhookTypeAppletTrigger,
		Data: make(chan []byte),
	}
}

// It removes a webhook from the Webhooks map
func RemoveWebhook(webhookName string) {
	delete(Webhooks, webhookName)
}

// It takes a webhook name, author name, and data, and pushes the data to the webhook's data channel
func WriteToWebhook(authorName string, webhookName string, data []byte) error {
	_, ok := Webhooks[webhookName]
	if !ok {
		return errors.New("Webhook does not exist")
	}

	History.PushFront(HistoryItem{
		ID:     strconv.Itoa(History.Len()),
		Author: authorName,
	})
	Webhooks[webhookName].Data <- data
	return nil
}

// GetWebhook returns a Webhook struct if the webhook mode matches the webhook name.
func GetWebhook(webhookName string, webhookMode WebhookMode) (Webhook, error) {
	if webhookMode != Webhooks[webhookName].Mode {
		return Webhook{}, errors.New("Webhook mode does not match")
	}
	return Webhooks[webhookName], nil
}
