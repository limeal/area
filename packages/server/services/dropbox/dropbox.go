package dropbox

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/dropbox/actions"
	"area-server/services/dropbox/reactions"
)

// It returns a static.Service object that describes the Dropbox service
func Descriptor() static.Service {
	return static.Service{
		Name:          "dropbox",
		Description:   "Dropbox is a file hosting service that offers cloud storage, file synchronization, personal cloud, and client software.",
		Authenticator: authenticators.GetAuthenticator("dropbox"),
		RateLimit:     2,
		Endpoints:     DropboxEndpoints(),
		Validators:    DropboxValidators(),
		Routes:        DropboxRoutes(),
		Actions: []static.ServiceArea{
			actions.DescriptorForDropboxActionNewFileInDirectory(),
			actions.DescriptorForDropboxActionUserShareAFileToNewEntity(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForDropboxReactionAddFileToDirectory(),
			reactions.DescriptorForDropboxReactionShareFileToUser(),
			reactions.DescriptorForDropboxReactionCreateOrUpdatePaper(),
		},
	}
}
