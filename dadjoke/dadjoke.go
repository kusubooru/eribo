package dadjoke

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://icanhazdadjoke.com/"
	defaultUserAgent = "Eribo (https://github.com/kusubooru/eribo)"
)

// Client is a client for the icanhazdadjoke API.
type Client struct {
	client *http.Client

	// User agent used when communicating with the icanhazdadjoke API.
	UserAgent string

	// Base URL for icanhazdadjoke API requests.
	BaseURL *url.URL
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot check response body: %v", err)
	}

	return fmt.Errorf("%v %v: %d %s",
		r.Request.Method, r.Request.URL,
		r.StatusCode, string(data))
}

// NewClient returns a new client for dadjokes.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		UserAgent: defaultUserAgent,
		BaseURL:   baseURL,
	}

	return c
}

// Random returns a random dadjoke.
func (c *Client) Random() (*Joke, error) {
	rel, err := url.Parse("/")
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	j := new(Joke)
	if err := json.NewDecoder(resp.Body).Decode(j); err != nil {
		return nil, err
	}
	return j, nil
}

// Joke represents a dadjoke.
type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// Random is a helper function that returns a random dadjoke.
func Random() (*Joke, error) {
	return NewClient(nil).Random()
}
