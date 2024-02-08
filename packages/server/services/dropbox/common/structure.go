package common

// Entry is a struct that represents a file or folder in Dropbox.
// @property {string} Tag - The type of file. This can be either file or folder.
// @property {string} Name - The name of the file or folder.
// @property {string} PathLower - The lowercased full path in the user's Dropbox. This always starts
// with a slash. This field will be null if the file or folder is not mounted.
// @property {string} PathDisplay - The path of the file or folder, including the file name.
// @property {int64} Size - The size of the file in bytes.
// @property {string} ClientModified - The last time the file was modified on Dropbox.
// @property {string} ServerModified - The last time the file was modified on Dropbox.
// @property {string} Rev - The revision of the file. This field is the same rev as elsewhere in the
// API and can be used to detect changes and avoid conflicts.
// @property {string} ID - A unique identifier for the file.
type Entry struct {
	Tag            string `json:".tag"`
	Name           string `json:"name"`
	PathLower      string `json:"path_lower"`
	PathDisplay    string `json:"path_display"`
	Size           int64  `json:"size"`
	ClientModified string `json:"client_modified"`
	ServerModified string `json:"server_modified"`
	Rev            string `json:"rev"`
	ID             string `json:"id"`
}
