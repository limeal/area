package github

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/github/actions"
	"area-server/services/github/reactions"
)

// It returns a static.Service object that describes the Github service
func Descriptor() static.Service {
	return static.Service{
		Name:          "github",
		Description:   "Github is a web-based hosting service for version control using Git. It is mostly used for computer code. It offers all of the distributed version control and source code management (SCM) functionality of Git as well as adding its own features.",
		Authenticator: authenticators.GetAuthenticator("github"),
		RateLimit:     3,
		Endpoints:     GithubEndpoints(),
		Validators:    GithubValidators(),
		Routes:        GithubRoutes(),
		Actions: []static.ServiceArea{
			actions.DescriptorForGithubActionAnyNewBranch(),
			actions.DescriptorForGithubActionAnyNewCommit(),
			actions.DescriptorForGithubActionAnyNewIssue(),
			actions.DescriptorForGithubActionAnyNewPullRequest(),
			actions.DescriptorForGithubActionAnyNewRelease(),
			actions.DescriptorForGithubActionNewRepository(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForGithubReactionAddCollaborator(),
			reactions.DescriptorForGithubReactionCreateNewGist(),
			reactions.DescriptorForGithubReactionCreateNewIssue(),
			reactions.DescriptorForGithubReactionCreateNewRelease(),
			reactions.DescriptorForGithubReactionCreateNewRepository(),
		},
	}
}
