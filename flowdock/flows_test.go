package flowdock_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	idOne     = "1"
	idTwo     = "2"
	idOrgFlow = "org:flow"
)

func (s *Suite) TestFlowsService_List() {
	s.mux.HandleFunc("/flows", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	flows, _, err := s.client.Flows.List(false, nil)
	s.NoError(err, "Flows.List returned error: %v", err)

	want := []flowdock.Flow{{ID: &idOne}, {ID: &idTwo}}
	s.Equal(want, flows, "Flows.List returned %+v, want %+v", flows, want)
}

func (s *Suite) TestFlowsService_List_all() {
	s.mux.HandleFunc("/flows/all", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	opt := flowdock.FlowsListOptions{User: true}
	flows, _, err := s.client.Flows.List(true, &opt)
	s.NoError(err, "Flows.List returned error: %v", err)

	want := []flowdock.Flow{{ID: &idOne}, {ID: &idTwo}}
	s.Equal(want, flows, "Flows.List returned %+v, want %+v", flows, want)
}

func (s *Suite) TestFlowsService_List_invalidOpt() {
	opt := new(flowdock.FlowsListOptions)
	_, _, err := s.client.Flows.List(true, opt)
	s.Error(err, "Flows.List expected an error")
}

func (s *Suite) TestFlowsService_Get() {
	s.mux.HandleFunc("/flows/orgname/flowname", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := s.client.Flows.Get("orgname", "flowname")
	s.NoError(err, "Flows.Get returned error: %v", err)

	want := &flowdock.Flow{ID: &idOne}
	s.Equal(want, flow, "Flows.Get returned %+v, want %+v", flow, want)
}

func (s *Suite) TestFlowsService_GetByID() {
	s.mux.HandleFunc("/flows/find", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		testFormValues(s.T(), r, values{"id": "orgname:flowname"})
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := s.client.Flows.GetByID("orgname:flowname")
	s.NoError(err, "Flows.Get returned error: %v", err)

	want := &flowdock.Flow{ID: &idOne}
	s.Equal(want, flow, "Flows.Get returned %+v, want %+v", flow, want)
}

func (s *Suite) TestFlowsService_Create() {
	s.mux.HandleFunc("/flows/org", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("POST", r.Method, "Request method = %v, want %v", r.Method, "POST")
		testFormValues(s.T(), r, values{"name": "flow"})
		fmt.Fprint(w, `{"id":"org:flow"}`)
	})

	opt := flowdock.FlowsCreateOptions{Name: "flow"}
	flow, _, err := s.client.Flows.Create("org", &opt)
	s.NoError(err, "Flows.Create returned error: %v", err)

	want := &flowdock.Flow{ID: &idOrgFlow}
	s.Equal(want, flow, "Flows.Create returned %+v, want %+v", flow, want)
}

func (s *Suite) TestFlowsService_Update() {
	truth := true
	input := &flowdock.Flow{Open: &truth}

	s.mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		v := new(flowdock.Flow)
		err := json.NewDecoder(r.Body).Decode(v)
		s.NoError(err)

		s.Equal("PUT", r.Method, "Request method = %v, want %v", r.Method, "PUT")
		s.Equal(input, v, "Request body = %+v, want %+v", v, input)
		fmt.Fprint(w, `{"id":"org:flow"}`)
	})

	flow, _, err := s.client.Flows.Update("org", "flow", input)
	s.NoError(err, "Flows.Update returned error: %v", err)

	want := &flowdock.Flow{ID: &idOrgFlow}
	s.Equal(want, flow, "Flows.Update returned %+v, want %+v", flow, want)
}
