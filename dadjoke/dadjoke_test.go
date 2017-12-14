package dadjoke

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	client *Client

	server *httptest.Server

	mux *http.ServeMux
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// icanhazdadjoke client configured to use test server
	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	// test default base URL
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient.BaseURL = %v, want %v", got, want)
	}

	// test default user agent
	if got, want := c.UserAgent, defaultUserAgent; got != want {
		t.Errorf("NewClient.UserAgent = %v, want %v", got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func testContentType(t *testing.T, r *http.Request, want string) {
	if got := r.Header.Get("Content-Type"); got != want {
		t.Errorf("Content-Type = %v, want %v", got, want)
	}
}

func testAccept(t *testing.T, r *http.Request, want string) {
	if got := r.Header.Get("Accept"); got != want {
		t.Errorf("Accept = %v, want %v", got, want)
	}
}

func testUserAgent(t *testing.T, r *http.Request, want string) {
	if got := r.Header.Get("User-Agent"); got != want {
		t.Errorf("User-Agent = %v, want %v", got, want)
	}
}

func TestClient_Random(t *testing.T) {
	setup()
	defer teardown()

	in := &Joke{
		ID:     "R7UfaahVfFd",
		Joke:   "My dog used to chase people on a bike a lot. It got so bad I had to take his bike away.",
		Status: 200,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testAccept(t, r, "application/json")
		testUserAgent(t, r, defaultUserAgent)
		if err := json.NewEncoder(w).Encode(in); err != nil {
			t.Fatal("could not encode test server response")
		}
	})

	joke, err := client.Random()
	if err != nil {
		t.Fatalf("client.Random() returned error %v", err)
	}
	if got, want := joke, in; !reflect.DeepEqual(got, want) {
		t.Errorf("client.Random() = %#v, want %#v", got, want)
	}
}
