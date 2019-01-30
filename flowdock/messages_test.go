package flowdock_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jtdoepke/go-flowdock/flowdock"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestMessagesService_Stream() {
	more := make(chan bool, 1)

	s.mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		testFormValues(s.T(), r, values{"access_token": "token"})
		w.Header().Set("Content-Type", "text/event-stream")

		var id int

		// send a message for each 'more'
		for {
			if !<-more {
				break
			}

			fmt.Fprintf(w, "id: %d\ndata: {\"event\":\"message\",\"content\":\"message %d\"}\n\n", id, id)
			w.(responseWriter).Flush()
			id++
		}
	})
	defer close(more)

	stream, _, err := s.client.Messages.Stream("token", "org", "flow")
	more <- true // tell test server to send a message

	s.NoError(err, "Messages.Stream returned error: %v", err)

	msg := <-stream

	s.Equal("message 0", msg.Content().String(), "expected message 0, got %v", msg.Content())
}

func (s *Suite) TestMessagesService_List() {
	var idOne = 3816534
	var eventOne = "message"
	var content = []string{"Hello NYC", "Hello World"}
	var idTwo = 45590
	var eventTwo = "message"

	s.mux.HandleFunc("/flows/org/flow/messages", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[
		  {
			"app":"chat",
			"sent":1317397485508,
			"uuid":"odHapx1VWp7WTrdQ",
			"tags":[],
			"flow": "deadbeefdeadbeef",
			"id":3816534,
			"event":"message",
			"content": "Hello NYC",
			"attachments": [],
			"user":"18"
		  },
		  {
			"app": "chat",
			"event": "message",
			"tags": [],
			"uuid": "4W_LQEybVaX-gJmi",
			"id": 45590,
			"flow": "deadbeefdeadbeef",
			"content": "Hello World",
			"sent": 1317715340213,
			"attachments": [],
			"user": "2"
		  }
		]`)
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	messages, _, err := s.client.Messages.List("org", "flow", nil)
	s.NoError(err, "Messages.List returned error: %v", err)

	want := []flowdock.Message{
		{
			ID:    &idOne,
			Event: &eventOne,
		},
		{
			ID:    &idTwo,
			Event: &eventTwo,
		},
	}

	for i, msg := range messages {
		s.Equal(*want[i].ID, *msg.ID, "Messages.List returned %+v, want %+v", *msg.ID, *want[i].ID)
		s.Equal(*want[i].Event, *msg.Event, "Messages.List returned %+v, want %+v", *msg.Event, *want[i].Event)
		s.Equal(content[i], msg.Content().String(), "Messages.List returned %+v, want %+v", msg.Content(), content[i])
	}
}

func (s *Suite) TestMessagesService_Create_message() {
	s.mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("POST", r.Method, "Request method = %v, want %v", r.Method, "POST")
		testFormValues(s.T(), r, values{"event": "message",
			"content": "Howdy-Doo @Jackie #awesome",
		})
		fmt.Fprint(w, `{
			"event": "message",
			"content": "Howdy-Doo @Jackie #awesome"
		}`)
	})

	opt := flowdock.MessagesCreateOptions{
		Event:   "message",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := s.client.Messages.Create(&opt)
	s.NoError(err, "Messages.Create returned error: %v", err)
	s.Equal(opt.Event, *message.Event, "Messages.Create returned %+v, want %+v", *message.Event, opt.Event)
	s.Equal(opt.Content, message.Content().String(), "Messages.Create returned %+v, want %+v", message.Content(), opt.Content)
}

func (s *Suite) TestMessagesService_Create_comment() {
	s.mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("POST", r.Method, "Request method = %v, want %v", r.Method, "POST")
		testFormValues(s.T(), r, values{"event": "comment",
			"content": "This is a comment",
		})
		fmt.Fprint(w, `{
			"event": "comment",
			"content":{ "title":"Title of parent", "text":"This is a comment" }
		}`)
	})

	opt := flowdock.MessagesCreateOptions{
		Event:   "comment",
		Content: "This is a comment",
	}
	message, _, err := s.client.Messages.CreateComment(&opt)
	s.NoError(err, "Messages.CreateComment returned error: %v", err)
	s.Equal(opt.Event, *message.Event, "Messages.Create returned %+v, want %+v", *message.Event, opt.Event)

	title := "Title of parent"
	text := "This is a comment"
	content := flowdock.CommentContent{Title: &title, Text: &text}
	messageContent := message.Content()
	s.Equal(&content, messageContent, "Messages.Create returned %+v, want %+v", messageContent, &content)
}

func TestCommentContent_String(t *testing.T) {
	t.Parallel()

	title := "Title of parent"
	text := "This is a comment"
	content := flowdock.CommentContent{Title: &title, Text: &text}

	want := "This is a comment"
	assert.Equal(t, want, *content.Text, "Messages.Create returned %+v, want %+v", *content.Text, want)
}
