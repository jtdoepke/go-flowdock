package flowdock_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	idOne     = "1"
	idTwo     = "2"
	idOrgFlow = "org:flow"
)

func TestFlowsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	flows, _, err := client.Flows.List(false, nil)
	assert.NoError(t, err, "Flows.List returned error: %v", err)

	want := []flowdock.Flow{{ID: &idOne}, {ID: &idTwo}}
	assert.Equal(t, want, flows, "Flows.List returned %+v, want %+v", flows, want)
}

func TestFlowsService_List_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	opt := flowdock.FlowsListOptions{User: true}
	flows, _, err := client.Flows.List(true, &opt)
	assert.NoError(t, err, "Flows.List returned error: %v", err)

	want := []flowdock.Flow{{ID: &idOne}, {ID: &idTwo}}
	assert.Equal(t, want, flows, "Flows.List returned %+v, want %+v", flows, want)
}

func TestFlowsService_List_invalidOpt(t *testing.T) {
	opt := new(flowdock.FlowsListOptions)
	_, _, err := client.Flows.List(true, opt)
	assert.Error(t, err, "Flows.List expected an error")
}

func TestFlowsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/orgname/flowname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := client.Flows.Get("orgname", "flowname")
	assert.NoError(t, err, "Flows.Get returned error: %v", err)

	want := &flowdock.Flow{ID: &idOne}
	assert.Equal(t, want, flow, "Flows.Get returned %+v, want %+v", flow, want)
}

func TestFlowsService_GetByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/find", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"id": "orgname:flowname"})
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := client.Flows.GetByID("orgname:flowname")
	assert.NoError(t, err, "Flows.Get returned error: %v", err)

	want := &flowdock.Flow{ID: &idOne}
	assert.Equal(t, want, flow, "Flows.Get returned %+v, want %+v", flow, want)
}

func TestFlowsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/org", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"name": "flow"})
		fmt.Fprint(w, `{"id":"org:flow"}`)
	})

	opt := flowdock.FlowsCreateOptions{Name: "flow"}
	flow, _, err := client.Flows.Create("org", &opt)
	assert.NoError(t, err, "Flows.Create returned error: %v", err)

	want := &flowdock.Flow{ID: &idOrgFlow}
	assert.Equal(t, want, flow, "Flows.Create returned %+v, want %+v", flow, want)
}

func TestFlowsService_Update(t *testing.T) {
	setup()
	defer teardown()

	truth := true
	input := &flowdock.Flow{Open: &truth}

	mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		v := new(flowdock.Flow)
		err := json.NewDecoder(r.Body).Decode(v)
		assert.NoError(t, err)

		testMethod(t, r, "PUT")
		assert.Equal(t, input, v, "Request body = %+v, want %+v", v, input)
		fmt.Fprint(w, `{"id":"org:flow"}`)
	})

	flow, _, err := client.Flows.Update("org", "flow", input)
	assert.NoError(t, err, "Flows.Update returned error: %v", err)

	want := &flowdock.Flow{ID: &idOrgFlow}
	assert.Equal(t, want, flow, "Flows.Update returned %+v, want %+v", flow, want)
}
