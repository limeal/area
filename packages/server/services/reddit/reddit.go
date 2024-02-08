package reddit

import (
	"area-server/authenticators"
	"area-server/classes/static"
	"area-server/services/reddit/actions"
	"area-server/services/reddit/reactions"
	"os"
)

// It returns a static.Service object that describes the service
func Descriptor() static.Service {

	_, p := os.LookupEnv("REDDIT_CLIENT_ID")
	if !p {
		panic("REDDIT_CLIENT_ID is not set")
	}
	_, pt := os.LookupEnv("REDDIT_SECRET_ID")
	if !pt {
		panic("REDDIT_SECRET_ID is not set")
	}

	return static.Service{
		Name:          "reddit",
		Description:   "Reddit is an American social news aggregation, web content rating, and discussion website.",
		Authenticator: authenticators.GetAuthenticator("reddit"),
		RateLimit:     4,
		Endpoints:     RedditEndpoints(),
		Validators:    RedditValidators(),
		Actions: []static.ServiceArea{
			actions.DescriptorForRedditActionNewCommentByYou(),
			actions.DescriptorForRedditActionNewDownvotedPostByYou(),
			actions.DescriptorForRedditActionNewPostByYou(),
			actions.DescriptorForRedditActionNewSavedPostByYou(),
			actions.DescriptorForRedditActionNewUpvotedPostByYou(),
		},
		Reactions: []static.ServiceArea{
			reactions.DescriptorForRedditReactionNewPostOnSubReddit(),
			reactions.DescriptorForRedditReactionSubmitNewLink(),
			reactions.DescriptorForRedditReactionSubmitNewPost(),
		},
	}
}
