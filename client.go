package tesla

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

// Provides the client and associated elements for interacting with the
// Tesla API
type Client struct {
	tokens oauth2.TokenSource
	HTTP   *http.Client
}

var (
	BaseURL      = "https://owner-api.teslamotors.com/api/1"
	ActiveClient *Client
)

// Generates a new client for the Tesla API
func NewClient(ts oauth2.TokenSource) (*Client, error) {
	client := &Client{
		tokens: ts,
		HTTP:   &http.Client{},
	}
	ActiveClient = client
	return client, nil
}

// // Calls an HTTP DELETE
func (c Client) delete(url string) error {
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err := c.processRequest(req)
	return err
}

// Calls an HTTP GET
func (c Client) get(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	return c.processRequest(req)
}

// Calls an HTTP POST with a JSON body
func (c Client) post(url string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	return c.processRequest(req)
}

// Calls an HTTP PUT
func (c Client) put(resource string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest("PUT", BaseURL+resource, bytes.NewBuffer(body))
	return c.processRequest(req)
}

// Processes a HTTP POST/PUT request
func (c Client) processRequest(req *http.Request) ([]byte, error) {
	if err := c.setHeaders(req); err != nil {
		return nil, err
	}
	res, err := c.HTTP.Do(req)
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
func (c Client) setHeaders(req *http.Request) error {
	token, err := c.tokens.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return nil
}
