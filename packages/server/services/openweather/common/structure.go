package common

// `Station` is a struct with fields `ID`, `CreatedAt`, `UpdatedAt`, `ExternalID`, `Name`, `Longitude`,
// `Latitude`, `Altitude`, and `Rank`.
//
// The `json` tags are used to tell the `json` package how to encode and decode the struct.
//
// The `json` package will look at the field names (`ID`, `CreatedAt`, etc.) and use those as the JSON
// keys.
//
// The `json` package will also look at the field types and encode/decode appropriately. For example
// @property {int} ID - The unique identifier for the station.
// @property {string} CreatedAt - The date and time the station was created.
// @property {string} UpdatedAt - The last time the station was updated.
// @property {string} ExternalID - The ID of the station in the external system.
// @property {string} Name - The name of the station
// @property {float64} Longitude - The longitude of the station.
// @property {float64} Latitude - The latitude of the station.
// @property {float64} Altitude - The altitude of the station in meters.
// @property {int} Rank - The rank of the station. This is used to determine the order in which the
// stations are displayed in the app.
type Station struct {
	ID         int     `json:"id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Altitude   float64 `json:"altitude"`
	Rank       int     `json:"rank"`
}
