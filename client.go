package tesla

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Required authorization credentials for the Tesla API
type Auth struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	URL          string
	StreamingURL string
	HTTPClient   HTTPDoer `json:"-"`
}

// The token and related elements returned after a successful auth
// by the Tesla API
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Expires     int64
}

// HTTPDoer is an http.Client implementation
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provides the client and associated elements for interacting with the
// Tesla API
type Client struct {
	Auth  *Auth
	Token *Token
	HTTP  HTTPDoer
}

var (
	AuthURL      = "https://owner-api.teslamotors.com/oauth/token"
	BaseURL      = "https://owner-api.teslamotors.com/api/1"
	ActiveClient *Client
)

// Generates a new client for the Tesla API
func NewClient(auth *Auth) (*Client, error) {
	if auth.URL == "" {
		auth.URL = BaseURL
	}
	if auth.StreamingURL == "" {
		auth.StreamingURL = StreamingURL
	}
	if auth.HTTPClient == nil {
		auth.HTTPClient = &http.Client{}
	}

	client := &Client{
		Auth: auth,
		HTTP: auth.HTTPClient,
	}
	token, err := client.authorize(auth)
	if err != nil {
		return nil, err
	}
	client.Token = token
	ActiveClient = client
	return client, nil
}

// NewClientWithToken Generates a new client for the Tesla API using an existing token
func NewClientWithToken(auth *Auth, token *Token) (*Client, error) {
	if auth.URL == "" {
		auth.URL = BaseURL
	}
	if auth.StreamingURL == "" {
		auth.StreamingURL = StreamingURL
	}

	client := &Client{
		Auth:  auth,
		HTTP:  &http.Client{},
		Token: token,
	}
	if client.TokenExpired() {
		return nil, errors.New("supplied token is expired")
	}
	ActiveClient = client
	return client, nil
}

// TokenExpired indicates whether an existing token is within an hour of expiration
func (c Client) TokenExpired() bool {
	exp := time.Unix(c.Token.Expires, 0)
	return time.Until(exp) < time.Duration(1*time.Hour)
}

// Authorizes against the Tesla API with the appropriate credentials
func (c Client) authorize(auth *Auth) (*Token, error) {
	now := time.Now()
	auth.GrantType = "password"
	data, _ := json.Marshal(auth)
	body, err := c.post(AuthURL, data)
	if err != nil {
		return nil, err
	}
	token := &Token{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return nil, err
	}
	token.Expires = now.Add(time.Second * time.Duration(token.ExpiresIn)).Unix()
	return token, nil
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
	c.setHeaders(req)
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
func (c Client) setHeaders(req *http.Request) {
	if c.Token != nil {
		req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}
