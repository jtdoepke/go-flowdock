package flowdock_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

func TestMessagesService_Stream(t *testing.T) {
	setup()
	defer teardown()
	more := make(chan bool, 1)

	mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"access_token": "token"})
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

	stream, _, err := client.Messages.Stream("token", "org", "flow")
	more <- true // tell test server to send a message

	assert.NoError(t, err, "Messages.Stream returned error: %v", err)

	msg := <-stream

	assert.Equal(t, "message 0", msg.Content().String(), "expected message 0, got %v", msg.Content())
}

func TestMessagesService_List(t *testing.T) {
	setup()
	defer teardown()
	var idOne = 3816534
	var eventOne = "message"
	var content = []string{"Hello NYC", "Hello World"}
	var idTwo = 45590
	var eventTwo = "message"

	mux.HandleFunc("/flows/org/flow/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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

	messages, _, err := client.Messages.List("org", "flow", nil)
	assert.NoError(t, err, "Messages.List returned error: %v", err)

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
		assert.Equal(t, *want[i].ID, *msg.ID, "Messages.List returned %+v, want %+v", *msg.ID, *want[i].ID)
		assert.Equal(t, *want[i].Event, *msg.Event, "Messages.List returned %+v, want %+v", *msg.Event, *want[i].Event)
		assert.Equal(t, content[i], msg.Content().String(), "Messages.List returned %+v, want %+v", msg.Content(), content[i])
	}
}

func TestMessagesService_Create_message(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"event": "message",
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
	message, _, err := client.Messages.Create(&opt)
	assert.NoError(t, err, "Messages.Create returned error: %v", err)
	assert.Equal(t, opt.Event, *message.Event, "Messages.Create returned %+v, want %+v", *message.Event, opt.Event)
	assert.Equal(t, opt.Content, message.Content().String(), "Messages.Create returned %+v, want %+v", message.Content(), opt.Content)
}

func TestMessagesService_Create_comment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"event": "comment",
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
	message, _, err := client.Messages.CreateComment(&opt)
	assert.NoError(t, err, "Messages.CreateComment returned error: %v", err)
	assert.Equal(t, opt.Event, *message.Event, "Messages.Create returned %+v, want %+v", *message.Event, opt.Event)

	title := "Title of parent"
	text := "This is a comment"
	content := flowdock.CommentContent{Title: &title, Text: &text}
	messageContent := message.Content()
	assert.Equal(t, &content, messageContent, "Messages.Create returned %+v, want %+v", messageContent, &content)
}

func TestCommentContent_String(t *testing.T) {
	title := "Title of parent"
	text := "This is a comment"
	content := flowdock.CommentContent{Title: &title, Text: &text}

	want := "This is a comment"
	assert.Equal(t, want, *content.Text, "Messages.Create returned %+v, want %+v", *content.Text, want)
}
