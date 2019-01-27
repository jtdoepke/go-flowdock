package flowdock_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jtdoepke/go-flowdock/flowdock"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Flowdock client being tested.
	client *flowdock.Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server

	// server is a test HTTP server used to provide mock API responses.
	streamServer *httptest.Server
)

// setup sets up a test HTTP server along with a flowdock.Client that is
// configured to talk to that test server. Tests should register handlers on mux
// which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	streamServer = httptest.NewServer(mux)

	// flowdock client configured to use test server
	client = flowdock.NewClient(nil)
	client.RestURL, _ = url.Parse(server.URL)
	client.StreamURL, _ = url.Parse(streamServer.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
	streamServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type responseWriter interface {
	http.ResponseWriter
	http.Flusher
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

func TestNewClient(t *testing.T) {
	c := flowdock.NewClient(nil)

	if c.RestURL.String() != flowdock.RestURL {
		t.Errorf("NewClient RestURL = %v, want %v", c.RestURL.String(), flowdock.RestURL)
	}
	if c.UserAgent != flowdock.UserAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, flowdock.UserAgent)
	}
}

func TestNewClientWithToken(t *testing.T) {
	token := "not-real-token"
	c := flowdock.NewClientWithToken(nil, token)

	url := fmt.Sprintf(flowdock.TokenRestURL, token)
	if c.RestURL.String() != url {
		t.Errorf("NewClientWithToken RestURL = %v, want %v", c.RestURL.String(), url)
	}
	if c.UserAgent != flowdock.UserAgent {
		t.Errorf("NewClientWithToken UserAgent = %v, want %v", c.UserAgent, flowdock.UserAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c := flowdock.NewClient(nil)

	name := "n"
	inURL, outURL := "/foo", flowdock.RestURL+"foo"
	inBody, outBody := &flowdock.Flow{Name: &name}, `{"name":"n"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, want %v", inBody, string(body), outBody)
	}

	// test that default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != flowdock.UserAgent {
		t.Errorf("NewRequest() User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := flowdock.NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	_, err := client.Do(req, body)
	if err != nil {
		t.Fatal(err)
	}

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m", 
                        "errors": [{"resource": "r", "field": "f", "code": "c"}]}`)),
	}
	err := flowdock.CheckResponse(res).(*flowdock.ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &flowdock.ErrorResponse{
		Response: res,
		Data: []byte(`{"message":"m", 
                        "errors": [{"resource": "r", "field": "f", "code": "c"}]}`),
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

// ensure that we properly handle API errors that do not contain a response
// body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := flowdock.CheckResponse(res).(*flowdock.ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &flowdock.ErrorResponse{
		Response: res,
		Data:     []byte{},
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := flowdock.ErrorResponse{Data: []byte("m"), Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}
