package flowdock_test

import (
	"fmt"
	"net/http"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

func (s *Suite) TestInboxService_Create() {
	s.mux.HandleFunc("/v1/messages/team_inbox/xxx", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("POST", r.Method, "Request method = %v, want %v", r.Method, "POST")
		testFormValues(s.T(), r, values{"subject": "a subject",
			"content": "Howdy-Doo @Jackie #awesome",
		})
		fmt.Fprint(w, `{}`)
	})

	opt := flowdock.InboxCreateOptions{
		Subject: "a subject",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := s.client.Inbox.Create("xxx", &opt)
	s.NoError(err, "Messages.Create returned error: %v", err)

	want := new(flowdock.Message)
	s.Equal(want, message, "Messages.Create returned \n%+v \nwant \n%+v", message, want)
}
