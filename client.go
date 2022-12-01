package tesla

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// Set to true to make the json decoder fail to decode if it encounters an
// unknown field.
const disallowUnknownFields = false

// OAuth2Config is the OAuth2 configuration for authenticating with the Tesla API.
var OAuth2Config = &oauth2.Config{
	ClientID:    "ownerapi",
	RedirectURL: "https://auth.tesla.com/void/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://auth.tesla.com/oauth2/v3/authorize",
		TokenURL:  "https://auth.tesla.com/oauth2/v3/token",
		AuthStyle: oauth2.AuthStyleInParams,
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
	ts           oauth2.TokenSource
	authHandler  *authHandler
}

// NewClient creates a new Tesla API client. You must provided one of WithToken or WithTokenFile
// functional options to initialize the client with an OAuth token.
func NewClient(ctx context.Context, options ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:      "https://owner-api.teslamotors.com/api/1",
		streamingURL: "https://streaming.vn.teslamotors.com",
		oc:           OAuth2Config,
	}

	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	// perform login if configured
	if client.authHandler != nil {
		if client.token != nil {
			return nil, errors.New("cannot have token and authorization options both")
		}

		var err error
		client.token, err = client.authHandler.login(ctx, client.oc)
		if err != nil {
			return nil, err
		}

		// wipe credentials
		client.authHandler = nil
	}

	if client.token == nil {
		return nil, errors.New("an OAuth2 token must be provided")
	}

	client.ts = client.oc.TokenSource(ctx, client.token)
	client.hc = oauth2.NewClient(ctx, client.ts)

	// use the Tesla UA transport
	client.hc.Transport = &Transport{RoundTripper: client.hc.Transport}

	return client, nil
}

// Token returns the oauth token
func (c Client) Token() (*oauth2.Token, error) {
	return c.ts.Token()
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
	decoder := json.NewDecoder(bytes.NewReader(body))
	if disallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err = decoder.Decode(out); err != nil {
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
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	return io.ReadAll(res.Body)
}

// Sets the required headers for calls to the Tesla API
func (c Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}
