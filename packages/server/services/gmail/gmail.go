package gmail

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/gmail/actions"
	"area-server/services/gmail/reactions"
)

// It returns a static.Service object that describes the Gmail service
func Descriptor() static.Service {
	return static.Service{
		Name:          "gmail",
		Description:   "Gmail is a free, advertising-supported email service developed by Google. Users can access Gmail on the web and using third-party programs that synchronize email content through POP or IMAP protocols. Gmail started as a limited beta release on April 1, 2004 and ended its testing phase on July 7, 2009.",
		Authenticator: authenticators.GetAuthenticator("google"),
		RateLimit:     2,
		More: &static.More{
			Avatar: true,
			Color:  "#EA4335",
		},
		Endpoints:  GmailEndpoints(),
		Validators: GmailValidators(),
		Actions: []static.ServiceArea{
			actions.DescriptorForGmailActionNewDraftMailCreated(),
			actions.DescriptorForGmailActionNewFilterCreated(),
			actions.DescriptorForGmailActionNewLabelCreated(),
			actions.DescriptorForGmailActionNewMailReceived(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForGmailReactionSendMail(),
			reactions.DescriptorForGmailReactionUpdateVacationMessage(),
		},
	}
}
