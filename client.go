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

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://auth.tesla.com/oauth2/v3/authorize",
	TokenURL: "https://auth.tesla.com/oauth2/v3/token",
}

const BaseURL = "https://owner-api.teslamotors.com/api/1"

// Provides the client and associated elements for interacting with the
// Tesla API
type Client struct {
	BaseURL      string
	StreamingURL string
	hc           *http.Client
}

// Generates a new client for the Tesla API
func NewClient(ctx context.Context, tok *oauth2.Token) (*Client, error) {
	config := &oauth2.Config{
		ClientID:    "ownerapi",
		RedirectURL: "https://auth.tesla.com/void/callback",
		Endpoint:    Endpoint,
		Scopes:      []string{"openid", "email", "offline_access"},
	}

	client := &Client{
		BaseURL:      BaseURL,
		StreamingURL: StreamingURL,
		hc:           config.Client(ctx, tok),
	}
	return client, nil
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
