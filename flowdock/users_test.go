package flowdock_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	userID1 int = 1
	userID2 int = 2
)

func TestUsersService_All(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := client.Users.All()
	if err != nil {
		t.Errorf("Users.All returned error: %v", err)
	}

	want := []flowdock.User{{ID: &userID1}, {ID: &userID2}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.All returned %+v, want %+v", users, want)
	}
}

func TestUsersService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/orgname/flowname/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := client.Users.List("orgname", "flowname")
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	want := []flowdock.User{{ID: &userID1}, {ID: &userID2}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.List returned %+v, want %+v", users, want)
	}
}

func TestUsersService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	user, _, err := client.Users.Get(userID1)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := flowdock.User{ID: &userID1}
	if !reflect.DeepEqual(user.ID, want.ID) {
		t.Errorf("Users.Get returned %+v, want %+v", user.ID, want.ID)
	}
}

func TestUsersService_Update(t *testing.T) {
	setup()
	defer teardown()

	nick := "new-nick"

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1, "nick":"new-nick"}`)
	})

	opts := &flowdock.UserUpdateOptions{
		Nick: "new-nick",
	}
	user, _, err := client.Users.Update(userID1, opts)
	if err != nil {
		t.Errorf("Users.Update returned error: %v", err)
	}

	want := flowdock.User{Nick: &nick}
	if !reflect.DeepEqual(user.Nick, want.Nick) {
		t.Errorf("Users.Update returned %+v, want %+v", user.Nick, want.Nick)
	}
}
