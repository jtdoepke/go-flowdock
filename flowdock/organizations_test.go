package flowdock_test

import (
	"fmt"
	"net/http"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	organizationID1 int = 1
	organizationID2 int = 2
)

func (s *Suite) TestOrganizationsService_All() {
	s.mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	organizations, _, err := s.client.Organizations.All()
	s.NoError(err, "Organizations.All returned error: %v", err)

	want := []flowdock.Organization{{ID: &organizationID1}, {ID: &organizationID2}}
	s.Equal(want, organizations, "Organizations.All returned %+v, want %+v", organizations, want)
}

func (s *Suite) TestOrganizationsService_GetByParameterizedName() {
	name := "parameterizedorgname"

	s.mux.HandleFunc("/organizations/parameterizedorgname", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `{"parameterized_name":"parameterizedorgname"}`)
	})

	organization, _, err := s.client.Organizations.GetByParameterizedName(name)
	s.NoError(err, "Organizations.GetByParameterizedName returned error: %v", err)

	want := flowdock.Organization{ParameterizedName: &name}
	s.Equal(
		want.ParameterizedName, organization.ParameterizedName,
		"Organizations.GetByParameterizedName returned %+v, want %+v", organization.ParameterizedName, want.ParameterizedName,
	)
}

func (s *Suite) TestOrganizationsService_GetByID() {
	s.mux.HandleFunc("/organizations/find?id=1", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("GET", r.Method, "Request method = %v, want %v", r.Method, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	organization, _, err := s.client.Organizations.GetByID(organizationID1)
	s.NoError(err, "Organizations.GetByID returned error: %v", err)

	want := flowdock.Organization{ID: &organizationID1}
	s.Equal(want.ID, organization.ID, "Organizations.GetByID returned %+v, want %+v", organization.ID, want.ID)
}

func (s *Suite) TestOrganizationsService_Update() {
	name := "new-name"

	s.mux.HandleFunc("/organizations/1", func(w http.ResponseWriter, r *http.Request) {
		s.Equal("PUT", r.Method, "Request method = %v, want %v", r.Method, "PUT")
		fmt.Fprint(w, `{"id":1, "name":"new-name"}`)
	})

	opts := &flowdock.OrganizationUpdateOptions{
		Name: name,
	}
	organization, _, err := s.client.Organizations.Update(organizationID1, opts)
	s.NoError(err, "Organizations.Update returned error: %v", err)

	want := flowdock.Organization{Name: &name}
	s.Equal(want.Name, organization.Name, "Organizations.Update returned %+v, want %+v", organization.Name, want.Name)
}
