package tesla

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/oauth2"
)

// github.com/uhthomas/tesla
func state() string {
	var b [9]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b[:])
}

// https://www.oauth.com/oauth2-servers/pkce/
func pkce() (verifier, challenge string, err error) {
	var p [87]byte
	if _, err := io.ReadFull(rand.Reader, p[:]); err != nil {
		return "", "", fmt.Errorf("rand read full: %w", err)
	}
	verifier = base64.RawURLEncoding.EncodeToString(p[:])
	b := sha256.Sum256([]byte(challenge))
	challenge = base64.RawURLEncoding.EncodeToString(b[:])
	return verifier, challenge, nil
}

type authHandler struct {
	auth     *auth
	username string
	password string
}

func (c *authHandler) login(ctx context.Context, oc *oauth2.Config) (*oauth2.Token, error) {
	verifier, challenge, err := pkce()
	if err != nil {
		return nil, err
	}

	c.auth.AuthURL = oc.AuthCodeURL(state(), oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	code, err := c.auth.Do(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	token, err := oc.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)

	return token, err
}

func defaultHandler() *authHandler {
	return &authHandler{
		auth: &auth{
			SelectDevice: mfaUnsupported,
			SolveCaptcha: captchaUnsupported,
		},
	}
}

// WithMFAHandler allows a consumer to provide a different configuration from the default.
func WithMFAHandler(handler func(context.Context, []Device) (Device, string, error)) ClientOption {
	return func(c *Client) error {
		if c.authHandler == nil {
			c.authHandler = defaultHandler()
		}

		c.authHandler.auth.SelectDevice = handler
		return nil
	}
}

func mfaUnsupported(_ context.Context, _ []Device) (Device, string, error) {
	return Device{}, "", errors.New("multi factor authentication is not supported")
}

// WithCaptchaHandler allows a consumer to provide a different configuration from the default.
func WithCaptchaHandler(handler func(context.Context, io.Reader) (string, error)) ClientOption {
	return func(c *Client) error {
		if c.authHandler == nil {
			c.authHandler = defaultHandler()
		}

		c.authHandler.auth.SolveCaptcha = handler
		return nil
	}
}

func captchaUnsupported(_ context.Context, _ io.Reader) (string, error) {
	return "", errors.New("captcha solving is not supported")
}

// WithCredentials allows a consumer to provide a different configuration from the default.
func WithCredentials(username, password string) ClientOption {
	return func(c *Client) error {
		if c.authHandler == nil {
			c.authHandler = defaultHandler()
		}

		c.authHandler.username = username
		c.authHandler.password = password
		return nil
	}
}
