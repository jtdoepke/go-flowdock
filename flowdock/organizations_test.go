package flowdock_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	organizationID1 int = 1
	organizationID2 int = 2
)

func TestOrganizationsService_All(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	organizations, _, err := client.Organizations.All()
	if err != nil {
		t.Errorf("Organizations.All returned error: %v", err)
	}

	want := []flowdock.Organization{{ID: &organizationID1}, {ID: &organizationID2}}
	if !reflect.DeepEqual(organizations, want) {
		t.Errorf("Organizations.All returned %+v, want %+v", organizations, want)
	}
}

func TestOrganizationsService_GetByParameterizedName(t *testing.T) {
	setup()
	defer teardown()

	name := "parameterizedorgname"

	mux.HandleFunc("/organizations/parameterizedorgname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"parameterized_name":"parameterizedorgname"}`)
	})

	organization, _, err := client.Organizations.GetByParameterizedName(name)
	if err != nil {
		t.Errorf("Organizations.GetByParameterizedName returned error: %v", err)
	}

	want := flowdock.Organization{ParameterizedName: &name}
	if !reflect.DeepEqual(organization.ParameterizedName, want.ParameterizedName) {
		t.Errorf("Organizations.GetByParameterizedName returned %+v, want %+v", organization.ParameterizedName, want.ParameterizedName)
	}
}

func TestOrganizationsService_GetByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/find?id=1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	organization, _, err := client.Organizations.GetByID(organizationID1)
	if err != nil {
		t.Errorf("Organizations.GetByID returned error: %v", err)
	}

	want := flowdock.Organization{ID: &organizationID1}
	if !reflect.DeepEqual(organization.ID, want.ID) {
		t.Errorf("Organizations.GetByID returned %+v, want %+v", organization.ID, want.ID)
	}
}

func TestOrganizationsService_Update(t *testing.T) {
	setup()
	defer teardown()

	name := "new-name"

	mux.HandleFunc("/organizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1, "name":"new-name"}`)
	})

	opts := &flowdock.OrganizationUpdateOptions{
		Name: name,
	}
	organization, _, err := client.Organizations.Update(organizationID1, opts)
	if err != nil {
		t.Errorf("Organizations.Update returned error: %v", err)
	}

	want := flowdock.Organization{Name: &name}
	if !reflect.DeepEqual(organization.Name, want.Name) {
		t.Errorf("Organizations.Update returned %+v, want %+v", organization.Name, want.Name)
	}
}
