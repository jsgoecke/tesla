package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bogosj/tesla"
	"github.com/manifoldco/promptui"
)

const (
	mfaPasscodeLength = 6
)

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
			if len(s) != mfaPasscodeLength {
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

	client, err := tesla.NewClient(
		context.Background(),
		tesla.WithMFAHandler(selectDevice),
		tesla.WithCredentials(username, password),
	)
	if err != nil {
		return err
	}

	t := client.Token()

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
