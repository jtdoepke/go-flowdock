package flowdock_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

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
	assert.NoError(t, err, "Organizations.All returned error: %v", err)

	want := []flowdock.Organization{{ID: &organizationID1}, {ID: &organizationID2}}
	assert.Equal(t, want, organizations, "Organizations.All returned %+v, want %+v", organizations, want)
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
	assert.NoError(t, err, "Organizations.GetByParameterizedName returned error: %v", err)

	want := flowdock.Organization{ParameterizedName: &name}
	assert.Equal(t,
		want.ParameterizedName, organization.ParameterizedName,
		"Organizations.GetByParameterizedName returned %+v, want %+v", organization.ParameterizedName, want.ParameterizedName,
	)
}

func TestOrganizationsService_GetByID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/find?id=1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	organization, _, err := client.Organizations.GetByID(organizationID1)
	assert.NoError(t, err, "Organizations.GetByID returned error: %v", err)

	want := flowdock.Organization{ID: &organizationID1}
	assert.Equal(t, want.ID, organization.ID, "Organizations.GetByID returned %+v, want %+v", organization.ID, want.ID)
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
	assert.NoError(t, err, "Organizations.Update returned error: %v", err)

	want := flowdock.Organization{Name: &name}
	assert.Equal(t, want.Name, organization.Name, "Organizations.Update returned %+v, want %+v", organization.Name, want.Name)
}
