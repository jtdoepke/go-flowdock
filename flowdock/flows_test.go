package flowdock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFlowsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	flows, _, err := client.Flows.List(false, nil)
	if err != nil {
		t.Errorf("Flows.List returned error: %v", err)
	}

	want := []Flow{{Id: "1"}, {Id: "2"}}
	if !reflect.DeepEqual(flows, want) {
		t.Errorf("Flows.List returned %+v, want %+v", flows, want)
	}
}

func TestFlowsService_List_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	opt := FlowsListOptions{User: true}
	flows, _, err := client.Flows.List(true, &opt)
	if err != nil {
		t.Errorf("Flows.List returned error: %v", err)
	}

	want := []Flow{{Id: "1"}, {Id: "2"}}
	if !reflect.DeepEqual(flows, want) {
		t.Errorf("Flows.List returned %+v, want %+v", flows, want)
	}
}

func TestFlowsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/orgname/flowname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := client.Flows.Get("orgname", "flowname")
	if err != nil {
		t.Errorf("Flows.Get returned error: %v", err)
	}

	want := &Flow{Id: "1"}
	if !reflect.DeepEqual(flow, want) {
		t.Errorf("Flows.Get returned %+v, want %+v", flow, want)
	}
}

func TestFlowsService_GetById(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/find", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"id": "orgname:flowname"})
		fmt.Fprint(w, `{"id":"1"}`)
	})

	flow, _, err := client.Flows.GetById("orgname:flowname")
	if err != nil {
		t.Errorf("Flows.Get returned error: %v", err)
	}

	want := &Flow{Id: "1"}
	if !reflect.DeepEqual(flow, want) {
		t.Errorf("Flows.Get returned %+v, want %+v", flow, want)
	}
}

func TestFlowsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/orgname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"name": "flowname"})
		fmt.Fprint(w, `{"id":"orgname:flowname"}`)
	})

	opt := FlowsCreateOptions{Name: "flowname"}
	flow, _, err := client.Flows.Create("orgname", &opt)
	if err != nil {
		t.Errorf("Flows.Create returned error: %v", err)
	}

	want := &Flow{Id: "orgname:flowname"}
	if !reflect.DeepEqual(flow, want) {
		t.Errorf("Flows.Create returned %+v, want %+v", flow, want)
	}
}

func TestFlowsService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := &Flow{Open: true}

	mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		v := new(Flow)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"id":"org:flow"}`)
	})

	flow, _, err := client.Flows.Update("org", "flow", input)
	if err != nil {
		t.Errorf("Flows.Update returned error: %v", err)
	}

	want := &Flow{Id: "org:flow"}
	if !reflect.DeepEqual(flow, want) {
		t.Errorf("Flows.Update returned %+v, want %+v", flow, want)
	}
}