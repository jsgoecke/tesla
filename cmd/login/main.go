package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/uhthomas/tesla"
	"golang.org/x/oauth2"
)

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

func selectDevice(ctx context.Context, devices []tesla.Device) (d tesla.Device, passcode string, err error) {
	var i int
	if len(devices) > 1 {
		var err error
		i, _, err = (&promptui.Select{
			Label:   "Device",
			Items:   devices,
			Pointer: promptui.PipeCursor,
		}).Run()
		if err != nil {
			return tesla.Device{}, "", fmt.Errorf("select device: %w", err)
		}
	}
	d = devices[i]

	passcode, err = (&promptui.Prompt{
		Label:   "Passcode",
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			if len(s) != 6 {
				return errors.New("len(s) != 6")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return tesla.Device{}, "", err
	}
	return d, passcode, nil
}

func Main(ctx context.Context) error {
	username, err := (&promptui.Prompt{
		Label:   "Username",
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("len(s) == 0")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return err
	}

	password, err := (&promptui.Prompt{
		Label:   "Password",
		Mask:    ' ',
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("len(s) == 0")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return err
	}

	verifier, challenge, err := pkce()
	if err != nil {
		return fmt.Errorf("pkce: %w", err)
	}

	c := &oauth2.Config{
		ClientID:     "ownerapi",
		ClientSecret: "",
		RedirectURL:  "https://auth.tesla.com/void/callback",
		Scopes:       []string{"openid email offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.tesla.com/oauth2/v3/authorize",
			TokenURL: "https://auth.tesla.com/oauth2/v3/token",
		},
	}

	code, err := (&tesla.Auth{
		AuthURL: c.AuthCodeURL(state(), oauth2.AccessTypeOffline,
			oauth2.SetAuthURLParam("code_challenge", challenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		),
		SelectDevice: selectDevice,
	}).Do(ctx, username, password)
	if err != nil {
		return err
	}

	t, err := c.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)
	if err != nil {
		return fmt.Errorf("exchange: %w", err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(c)
	e.Encode(t)

	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
