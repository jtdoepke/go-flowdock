package flowdock_test

import (
	"fmt"
	"net/http"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	userID1 int = 1
	userID2 int = 2
)

func (s *Suite) TestUsersService_All() {
	s.mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := s.client.Users.All()
	s.NoError(err, "Users.All returned error: %v", err)

	want := []flowdock.User{{ID: &userID1}, {ID: &userID2}}
	s.Equal(want, users, "Users.All returned %+v, want %+v", users, want)
}

func (s *Suite) TestUsersService_List() {
	s.mux.HandleFunc("/flows/orgname/flowname/users", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := s.client.Users.List("orgname", "flowname")
	s.NoError(err, "Users.List returned error: %v", err)

	want := []flowdock.User{{ID: &userID1}, {ID: &userID2}}
	s.Equal(want, users, "Users.List returned %+v, want %+v", users, want)
}

func (s *Suite) TestUsersService_Get() {
	s.mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	user, _, err := s.client.Users.Get(userID1)
	s.NoError(err, "Users.Get returned error: %v", err)

	want := flowdock.User{ID: &userID1}
	s.Equal(want.ID, user.ID, "Users.Get returned %+v, want %+v", user.ID, want.ID)
}

func (s *Suite) TestUsersService_Update() {
	nick := "new-nick"

	s.mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("PUT", r.Method, "Request method = %v, want %v", r.Method, "PUT")
		fmt.Fprint(w, `{"id":1, "nick":"new-nick"}`)
	})

	opts := &flowdock.UserUpdateOptions{
		Nick: "new-nick",
	}
	user, _, err := s.client.Users.Update(userID1, opts)
	s.NoError(err, "Users.Update returned error: %v", err)

	want := flowdock.User{Nick: &nick}
	s.Equal(want.Nick, user.Nick, "Users.Update returned %+v, want %+v", user.Nick, want.Nick)
}
