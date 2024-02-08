package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// HelloData is a struct with an Op field of type int and a D field of type struct with a
// HeartbeatInterval field of type int.
// @property {int} Op - The operation code.
// @property D - This is the data that is sent with the event.
type HelloData struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

// Ping is a struct with two fields, Op and D, both of which are integers.
// @property {int} Op - The operation code. This is always 1.
// @property {int} D - The current Unix time in milliseconds.
type Ping struct {
	Op int `json:"op"`
	D  int `json:"d"`
}

// Activity is a struct with two fields, Name and Type, both of which are strings.
// @property {string} Name - The name of the activity.
// @property {int} Type - The type of activity. This can be one of the following:
type Activity struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// `IdentityPresence` is a struct that contains a slice of `Activity`s, a string, an int, and a bool.
// @property {[]Activity} Activities - An array of activities the user is doing.
// @property {string} Status - The user's current status.
// @property {int} Since - The Unix timestamp of when the user's presence was last updated.
// @property {bool} Afk - Whether or not the user is AFK
type IdentityPresence struct {
	Activities []Activity `json:"activities"`
	Status     string     `json:"status"`
	Since      int        `json:"since"`
	Afk        bool       `json:"afk"`
}

// IdentityProperties is a struct that contains the properties of an identity.
// @property {string} Os - The operating system of the device.
// @property {string} Browser - The browser name.
// @property {string} Device - The device type.
type IdentityProperties struct {
	Os      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

// `IdentifyD` is a struct with four fields, `Token`, `Intents`, `Properties`, and `Presence`.
//
// The first field, `Token`, is a string. The second field, `Intents`, is an int. The third field,
// `Properties`, is a struct of type `IdentityProperties`. The fourth field, `Presence`, is a struct of
// type `IdentityPresence`.
// @property {string} Token - The token of the bot that you want to identify.
// @property {int} Intents - This is a bitmask of the intents that the bot wishes to receive.
// @property {IdentityProperties} Properties - This is the object that contains the user's username,
// avatar, and discriminator.
// @property {IdentityPresence} Presence - The presence of the user.
type IdentifyD struct {
	Token      string             `json:"token"`
	Intents    int                `json:"intents"`
	Properties IdentityProperties `json:"properties"`
	Presence   IdentityPresence   `json:"presence"`
}

// `Identify` is a struct with two fields, `Op` and `D`. `Op` is an `int` and `D` is a `IdentifyD`
// struct.
// @property {int} Op - The operation code. This is always 2.
// @property {IdentifyD} D - The data for the identify payload.
type Identify struct {
	Op int       `json:"op"`
	D  IdentifyD `json:"d"`
}

// `Event` is a struct with four fields, `Op`, `T`, `S`, and `D`.
//
// The `Op` field is an integer, `T` is a string, `S` is an integer, and `D` is a map of strings to
// interfaces.
//
// The `Op` field is the operation code of the event.
//
// The `T` field is the event name.
//
// The `S` field is the event's sequence number.
//
// The `D` field is the event's data.
// @property {int} Op - The operation code.
// @property {string} T - The type of event.
// @property {int} S - Sequence number
// @property D - The data of the event.
type Event struct {
	Op int                    `json:"op"`
	T  string                 `json:"t"`
	S  int                    `json:"s"`
	D  map[string]interface{} `json:"d"`
}

// This is a method that is attached to the Event struct. It takes a string as an argument and
// returns an interface. It returns the value of the key in the D field of the Event struct.
func (e *Event) GetEvent(key string) interface{} {
	return e.D[key]
}

// `DiscordGateway` is a struct that contains a websocket connection, an integer, a channel of type
// bool, and a channel of type Event.
// @property WS - The websocket connection to the Discord gateway.
// @property {int} HeartbeatInterval - The interval at which the gateway will send a heartbeat to the
// Discord API.
// @property Interrupt - This is a channel that will be used to interrupt the heartbeat loop.
// @property EventChan - This is the channel that the gateway will send events to.
type DiscordGateway struct {
	WS                *websocket.Conn `json:"-"`
	HeartbeatInterval int             `json:"-"`
	Interrupt         chan bool       `json:"-"`
	EventChan         chan Event      `json:"-"`
}

// This is the function that connects to the Discord gateway. It uses the `websocket` package to
// connect to the gateway. It then receives the first message from the gateway and unmarshals it into a
// `HelloData` struct. It then sets the `HeartbeatInterval` field of the `DiscordGateway` struct to the
// `HeartbeatInterval` field of the `HelloData` struct.
func (g *DiscordGateway) Connect() error {
	c, _, err := websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=6&encoding=json", nil)

	if err != nil {
		return err
	}
	fmt.Println("Connected To Websocket!")
	g.WS = c
	msg, errr := g.Receive()
	if errr != nil {
		return errr
	}

	hello := HelloData{}
	err = json.Unmarshal(msg, &hello)
	if err != nil {
		return err
	}
	g.HeartbeatInterval = hello.D.HeartbeatInterval
	return nil
}

// Sending a message to the Discord gateway.
func (g *DiscordGateway) Send(msg []byte) error {
	err := g.WS.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

// Receiving a message from the Discord gateway.
func (g *DiscordGateway) Receive() ([]byte, error) {
	_, message, err := g.WS.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Creating a Ping struct and then marshalling it into a JSON string. It then sends the
// JSON string to the Discord gateway.
func (g *DiscordGateway) Ping() error {
	ping := Ping{
		Op: 1,
		D:  g.HeartbeatInterval,
	}
	msg, err := json.Marshal(ping)
	if err != nil {
		return err
	}
	err = g.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Creating an `Identify` struct and then marshalling it into a JSON string. It then sends
// the JSON string to the Discord gateway.
func (g *DiscordGateway) Identify() error {
	identify := Identify{
		Op: 2,
		D: IdentifyD{
			Token:   os.Getenv("DISCORD_BOT_TOKEN"),
			Intents: 131071,
			Properties: IdentityProperties{
				Os:      "linux",
				Browser: "chrome",
				Device:  "chrome",
			},
			Presence: IdentityPresence{
				Activities: []Activity{
					{
						Name: "Im watching you",
						Type: 0,
					},
				},
				Status: "dnd",
				Since:  91879201,
				Afk:    false,
			},
		},
	}

	msg, err := json.Marshal(identify)
	if err != nil {
		return err
	}

	err = g.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// This function is starting the Discord gateway. It is connecting to the gateway, sending a
// ping, receiving a message, identifying the bot, receiving a message, and then starting two
// goroutines.
// One goroutine is sending messages to the `message` channel and the other is sending pings to the
// Discord gateway. It then enters a for loop that will receive messages from the `message` channel and
// unmarshal them into an `Event` struct. It will then send the `Event` struct to the `EventChan`
// channel.
// If the `done` channel is closed, it will return nil. If the `Interrupt` channel is closed, it will
// write a close message to the Discord gateway and then return nil.
func (g *DiscordGateway) Start() error {
	err := g.Connect()
	if err != nil {
		return err
	}

	defer g.Close()
	done := make(chan struct{})

	// First ping
	err = g.Ping()
	if err != nil {
		return err
	}

	_, pingError := g.Receive()
	if pingError != nil {
		return pingError
	}

	err = g.Identify()
	if err != nil {
		return err
	}

	_, identityError := g.Receive()
	if identityError != nil {
		return identityError
	}

	message := make(chan []byte)

	go func() {
		defer close(done)
		for {
			time.Sleep(2 * time.Second)
			_, msg, err := g.WS.ReadMessage()
			if err != nil {
				return
			}
			message <- msg
		}
	}()

	go func() {
		defer close(done)
		for {
			time.Sleep(time.Duration(g.HeartbeatInterval) * time.Millisecond)
			err := g.Ping()
			if err != nil {
				return
			}
		}
	}()

	for {
		select {
		case msg := <-message:
			event := Event{}
			err := json.Unmarshal(msg, &event)
			if err != nil {
				return err
			}
			g.EventChan <- event
		case <-done:
			return nil
		case <-g.Interrupt:
			err := g.WS.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

// Closing the websocket connection.
func (g *DiscordGateway) Close() error {
	return g.WS.Close()
}

// Sending a message to the `Interrupt` channel.
func (g *DiscordGateway) Stop() error {
	g.Interrupt <- true
	return nil
}
