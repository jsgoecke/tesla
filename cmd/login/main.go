package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bogosj/tesla"
	"github.com/manifoldco/promptui"
	"golang.org/x/oauth2"
)

const (
	mfaPasscodeLength = 6
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

func selectDevice(ctx context.Context, devices []device) (d device, passcode string, err error) {
	var i int
	if len(devices) > 1 {
		var err error
		i, _, err = (&promptui.Select{
			Label:   "Device",
			Items:   devices,
			Pointer: promptui.PipeCursor,
		}).Run()
		if err != nil {
			return device{}, "", fmt.Errorf("select device: %w", err)
		}
	}
	d = devices[i]

	passcode, err = (&promptui.Prompt{
		Label:   "Passcode",
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			if len(s) != mfaPasscodeLength {
				return errors.New("len(s) != 6")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return device{}, "", err
	}
	return d, passcode, nil
}

func getUsernameAndPassword() (string, string, error) {
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
		return "", "", err
	}

	password, err := (&promptui.Prompt{
		Label:   "Password",
		Mask:    '*',
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("len(s) == 0")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func shortLongStringFlag(name, short, value, usage string) *string {
	s := flag.String(name, value, usage)
	flag.StringVar(s, short, value, usage)
	return s
}

func login(ctx context.Context) error {
	out := shortLongStringFlag("out", "o", "", "Token JSON output path. Leave blank or use '-' to write to stdout.")
	flag.Parse()

	username, password, err := getUsernameAndPassword()
	if err != nil {
		log.Fatal(err)
	}

	verifier, challenge, err := pkce()
	if err != nil {
		return fmt.Errorf("pkce: %w", err)
	}

	c := tesla.DefaultOAuth2Config

	code, err := (&auth{
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

	w := os.Stdout
	switch out := *out; out {
	case "", "-":
	default:
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir all: %w", err)
		}
		f, err := os.OpenFile(filepath.Clean(out), os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()
		w = f
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	if err := e.Encode(t); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}

func main() {
	if err := login(context.Background()); err != nil {
		log.Fatal(err)
	}
}
