package flowdock

import (
	"fmt"
)

// Content should be implemented by any value that is parsed into
// Message.RawContent. Its API will likly expand as more Message types are
// implemented.
type Content interface {
	String() string
}

// MessageContent represents a Message's Content when Message.Event is "message"
type MessageContent string

// Return the string version of a MessageContent
//
func (c *MessageContent) String() string {
	return string(*c)
}

// JsonContent is the default type for Message.Content() that does not have its
// own explicit type.
//
type JsonContent string

// Unmarshal the json data into JsonContent (i.e. just a string really)
//
// This just casts a byte data into a JsonContent
func (c *JsonContent) UnmarshalJSON(data []byte) error {
	*c = JsonContent(data)
	return nil
}

// Return the string version of a JsonContent
//
func (c *JsonContent) String() string {
	return string(*c)
}

// CommentContent represents a Message's Content when Message.Event is "comment"
type CommentContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
}

// Return the string version of a CommentContent
//
// It returns the *CommentContent.Text
func (c *CommentContent) String() string {
	return *c.Text
}

// VCS (i.e. Github)
type VcsContent struct {
	Issue struct {
		URL *string `json:"url"`
	} `json:"issue,omitempty"`
	PullRequest struct {
		URL *string `json:"url"`
	} `json:"pull_request,omitempty"`
	Pusher struct {
		Name *string `json:"name"`
	} `json:"pusher,omitempty"`
	Sender struct {
		Login *string `json:"login"`
	} `json:"sender,omitempty"`
	Repository struct {
		Name *string `json:"name"`
	} `json:"repository,omitempty"`
	Event      *string `json:"event,omitempty"`
	CompareURL *string `json:"compare,omitempty"`
}

// Return the string version of a VcsContent
//
// It returns the *CommentContent.Text
func (c *VcsContent) String() string {
	var user, url string
	event := *c.Event
	name := *c.Repository.Name

	if c.Pusher.Name != nil {
		user = *c.Pusher.Name
	} else if c.Sender.Login != nil {
		user = *c.Sender.Login
	} else {
		user = "Unknown"
	}

	if c.CompareURL != nil {
		url = *c.CompareURL
	} else if c.PullRequest.URL != nil {
		url = *c.PullRequest.URL
	} else if c.Issue.URL != nil {
		url = *c.Issue.URL
	}

	return fmt.Sprintf("%s: %s by %s %s", name, event, user, url)
}
