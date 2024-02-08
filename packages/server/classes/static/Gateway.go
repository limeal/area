package static

// Gateway is an interface that defines the methods that a gateway must implement.
// @property {error} Connect - This is the method that will be called to connect to the gateway.
// @property {error} Send - This is the method that will be used to send messages to the gateway.
// @property Receive - This is the channel that will receive messages from the gateway.
// @property {error} Ping - This is a method that will be called periodically to check if the
// connection is still alive.
// @property {error} Start - Starts the gateway.
// @property {error} Close - Closes the connection to the gateway.
// @property {error} Stop - This is a channel that is closed when the gateway is stopped.
type Gateway interface {
	Connect() error
	Send(msg []byte) error
	Receive() ([]byte, error)
	Ping() error
	Start() error
	Close() error
	Stop() error
}
