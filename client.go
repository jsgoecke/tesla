package tesla

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

// DefaultOAuth2Config is the OAuth2 configuration for authenticating with the Tesla API.
var DefaultOAuth2Config = &oauth2.Config{
	ClientID:    "ownerapi",
	RedirectURL: "https://auth.tesla.com/void/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://auth.tesla.com/oauth2/v3/authorize",
		TokenURL: "https://auth.tesla.com/oauth2/v3/token",
	},
	Scopes: []string{"openid", "email", "offline_access"},
}

// Client provides the client and associated elements for interacting with the Tesla API.
type Client struct {
	baseURL      string
	streamingURL string
	hc           *http.Client
	oc           *oauth2.Config
	token        *oauth2.Token
}

type ClientOption func(c *Client) error

// WithToken provides an oauth2.Token to the client for auth.
func WithToken(t *oauth2.Token) ClientOption {
	return func(c *Client) error {
		c.token = t
		return nil
	}
}

// WithTokenFile reads a JSON serialized oauth2.Token struct from disk and provides it
// to the client for auth.
func WithTokenFile(path string) ClientOption {
	t, err := loadToken(path)
	if err != nil {
		return func(c *Client) error {
			return err
		}
	}
	return WithToken(t)
}

// WithBaseURL provides a method to set the base URL for standard API calls to differ
// from the default.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}

// WithStreamingURL provides a method to set the base URL for streaming API calls to differ
// from the default.
func WithStreamingURL(url string) ClientOption {
	return func(c *Client) error {
		c.streamingURL = url
		return nil
	}
}

// WithOAuth2Config allows a consumer to provide a different configuration from the default.
func WithOAuth2Config(oc *oauth2.Config) ClientOption {
	return func(c *Client) error {
		c.oc = oc
		return nil
	}
}

// New creates a new Tesla API client. You must provided one of WithToken or WithTokenFile
// functional options to initialize the client with an OAuth token.
func NewClient(ctx context.Context, options ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:      "https://owner-api.teslamotors.com/api/1",
		streamingURL: "https://streaming.vn.teslamotors.com",
		oc:           DefaultOAuth2Config,
	}

	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	if client.token == nil {
		return nil, errors.New("an OAuth2 token must be provided")
	}

	client.hc = client.oc.Client(ctx, client.token)

	return client, nil
}

func loadToken(path string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tok := new(oauth2.Token)
	if err := json.Unmarshal(b, tok); err != nil {
		return nil, err
	}
	return tok, nil
}

// Calls an HTTP GET
func (c Client) get(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	return c.processRequest(req)
}

// getJSON performs an HTTP GET and then unmarshals the result into the provided struct.
func (c Client) getJSON(url string, out interface{}) error {
	body, err := c.get(url)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, out); err != nil {
		return err
	}
	return nil
}

// Calls an HTTP POST with a JSON body
func (c Client) post(url string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	return c.processRequest(req)
}

// Processes a HTTP POST/PUT request
func (c Client) processRequest(req *http.Request) ([]byte, error) {
	c.setHeaders(req)
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Sets the required headers for calls to the Tesla API
func (c Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}
