package shared

// AreaResponse is a struct that contains an error, a boolean, and a map of strings to interfaces.
// @property {error} Error - This is the error that is returned from the server.
// @property {bool} Success - This is a boolean value that indicates whether the request was successful
// or not.
// @property Data - This is the data that will be returned to the client.
type AreaResponse struct {
	Error   error                  `json:"error"`
	Success bool                   `json:"sucess"`
	Data    map[string]interface{} `json:"data"`
}
