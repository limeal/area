package common

// PartialCommit is a struct with two fields, Sha and Url, both of which are strings.
// @property {string} Sha - The SHA of the commit.
// @property {string} Url - The URL of the commit.
type PartialCommit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

// `Branch` is a struct with three fields, `Name` of type `string`, `Commit` of type `PartialCommit`,
// and `Protected` of type `bool`.
// @property {string} Name - The name of the branch.
// @property {PartialCommit} Commit - The commit that the branch points to.
// @property {bool} Protected - Whether or not the branch is protected.
type Branch struct {
	Name      string        `json:"name"`
	Commit    PartialCommit `json:"commit"`
	Protected bool          `json:"protected"`
}

// `PartialAuthor` is a struct with two fields, `Name` and `Email`, both of which are strings.
//
// The `json` tags are used to tell the `encoding/json` package how to marshal and unmarshal the
// struct.
//
// The `json` tags are also used by the `go-swagger` tool to generate the API documentation.
// @property {string} Name - The name of the author.
// @property {string} Email - The email address of the author.
type PartialAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// `Commit` is a struct with a `Sha` field of type `string` and an `ICommit` field of type `struct`
// with a `Message` field of type `string` and an `Author` field of type `PartialAuthor`.
// @property {string} Sha - The SHA of the commit.
// @property ICommit - The commit object
type Commit struct {
	Sha     string `json:"sha"`
	ICommit struct {
		Message string        `json:"message"`
		Author  PartialAuthor `json:"author"`
	} `json:"commit"`
}

// `PartialUser` is a struct with a single field, `Login`, which is a string.
// @property {string} Login - The username of the user
type PartialUser struct {
	Login string `json:"login"`
}

// A Label is a struct with four fields: ID, Name, Color, and Description.
// @property {int} ID - The ID of the label.
// @property {string} Name - The name of the label.
// @property {string} Color - The color of the label, in 6-digit hex notation.
// @property {string} Description - The description of the label.
type Label struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// `Issue` is a struct that contains a bunch of fields, some of which are strings, some of which are
// integers, and some of which are other structs.
// @property {int} ID - The ID of the issue.
// @property {int} Number - The issue number.
// @property {string} Title - The title of the issue.
// @property {string} Body - The contents of the issue.
// @property {string} State - The state of the issue. Either open or closed.
// @property {string} HTMLURL - The URL to the issue on GitHub.
// @property {string} CreatedAt - The date and time the issue was created.
// @property {PartialUser} User - The user who created the issue.
// @property {[]Label} Labels - The labels that are attached to the issue.
// @property {[]PartialUser} Assignees - The list of users assigned to this issue.
// @property {Repository} Repository - The repository that the issue belongs to.
type Issue struct {
	ID         int           `json:"id"`
	Number     int           `json:"number"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	State      string        `json:"state"`
	HTMLURL    string        `json:"html_url"`
	CreatedAt  string        `json:"created_at"`
	User       PartialUser   `json:"user"`
	Labels     []Label       `json:"labels"`
	Assignees  []PartialUser `json:"assignees"`
	Repository Repository    `json:"repository"`
}

// `Release` is a struct with fields `TagName`, `Name`, `Body`, `ID`, `HTMLURL`, `Author`, and
// `TarBallURL`.
//
// The `json` tags are used to tell the `encoding/json` package how to marshal and unmarshal the
// struct.
//
// The `PartialUser` type is a struct with fields `Login` and `ID`.
//
// The `Release` type is a struct with fields `TagName`, `Name`, `Body`, `ID`, `HTMLURL`, `Author`, and
// `Tar
// @property {string} TagName - The name of the tag.
// @property {string} Name - The name of the release.
// @property {string} Body - The body of the release.
// @property {int} ID - The ID of the release.
// @property {string} HTMLURL - The URL to the release on GitHub.
// @property {PartialUser} Author - The user who created the release.
// @property {string} TarBallURL - The URL to download the release as a tarball.
type Release struct {
	TagName    string      `json:"tag_name"`
	Name       string      `json:"name"`
	Body       string      `json:"body"`
	ID         int         `json:"id"`
	HTMLURL    string      `json:"html_url"`
	Author     PartialUser `json:"author"`
	TarBallURL string      `json:"tarball_url"`
}

// A Collaborator is a type that has a field called Login that is a string.
// @property {string} Login - The username of the collaborator.
type Collaborator struct {
	Login string `json:"login"`
}

// A PullRequest is a struct with a bunch of fields.
// @property {int} ID - The ID of the pull request.
// @property {int} Number - The number of the pull request.
// @property {string} Title - The title of the pull request.
// @property {string} Body - The body of the pull request.
// @property {string} State - The state of the pull request. Can be either open or closed.
// @property {string} HTMLURL - The URL to the pull request on GitHub.
// @property {PartialUser} User - The user who created the pull request.
// @property Head - The branch where your changes are implemented.
type PullRequest struct {
	ID      int         `json:"id"`
	Number  int         `json:"number"`
	Title   string      `json:"title"`
	Body    string      `json:"body"`
	State   string      `json:"state"`
	HTMLURL string      `json:"html_url"`
	User    PartialUser `json:"user"`
	Head    struct {
		Ref  string     `json:"ref"`
		Repo Repository `json:"repo"`
	} `json:"head"`
}

// Repository is a struct that contains the fields ID, Name, Description, HTMLURL, Owner, Private, and
// CreatedAt.
// @property {int} ID - The repository's ID.
// @property {string} Name - The name of the repository.
// @property {string} Description - The description of the repository.
// @property {string} HTMLURL - The URL to the repository on GitHub.
// @property {PartialUser} Owner - The user who owns the repository.
// @property {bool} Private - Whether the repository is private or not.
// @property {string} CreatedAt - The date and time when the repository was created.
type Repository struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	HTMLURL     string      `json:"html_url"`
	Owner       PartialUser `json:"owner"`
	Private     bool        `json:"private"`
	CreatedAt   string      `json:"created_at"`
}

// `GistFile` is a struct that has a single field, `Content`, which is a string.
//
// Now, let's create a new `GistFile` and assign it to a variable named `file`.
// @property {string} Content - The content of the file.
type GistFile struct {
	Content string `json:"content"`
}

// `CreateGist` is a struct that contains a string, a bool, and a map of strings to `GistFile`s.
// @property {string} Description - The description of the gist.
// @property {bool} Public - This is a boolean value that determines whether the gist is public or
// private.
// @property Files - A map of the files in the gist. The key of the map is the filename and the value
// is the GistFile struct.
type CreateGist struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

// CreateIssue is a struct that contains the fields Title, Body, Assignees, Labels, and Milestone.
// @property {string} Title - The title of the issue.
// @property {string} Body - The contents of the issue.
// @property {[]string} Assignees - An array of user logins to assign to this issue.
// @property {[]string} Labels - The labels to associate with this issue. Pass one or more labels to
// replace the set of labels on this issue. Send an empty array ([]) to clear all labels from the
// issue.
// @property {int} Milestone - The milestone to associate this issue with.
type CreateIssue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Labels    []string `json:"labels"`
	Milestone int      `json:"milestone"`
}

// `CreateRelease` is a struct that contains the fields `TagName`, `TargetCommitish`, `Name`, `Body`,
// `Draft`, and `Prerelease`.
// @property {string} TagName - The name of the tag.
// @property {string} TargetCommitish - The branch name to create the release from.
// @property {string} Name - The name of the release.
// @property {string} Body - The text describing the release.
// @property {bool} Draft - true if the release is a draft, false if it's a published release.
// @property {bool} Prerelease - true if the release is a prerelease, false if it's a full release.
type CreateRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
}
