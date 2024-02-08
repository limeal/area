package config

type ServerMode int

// Creating an enum.
const (
	Session ServerMode = iota
	Token
)

// `Config` is a struct that contains a `ServerMode` (which is an enum), an `int`, and a `bool`.
// @property {ServerMode} Mode - This is the mode of the server. It can be either "dev" or "prod".
// @property {int} TokenDuration - The duration of the token in seconds.
// @property {bool} HTTPS - If true, the server will run on HTTPS.
type Config struct {
	Mode          ServerMode
	TokenDuration int
	HTTPS         bool
}

// Creating a global variable called CFG that is a pointer to a Config struct.
var CFG = &Config{
	Mode:          Token,
	TokenDuration: 60 * 60 * 24 * 7,
	HTTPS:         false,
}
