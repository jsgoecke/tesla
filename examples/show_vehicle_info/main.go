package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bogosj/tesla"
	"golang.org/x/oauth2"
)

var tokenPath = flag.String("token", "", "path to token file")

func main() {
	flag.Parse()

	if *tokenPath == "" {
		fmt.Println("--token must be specified")
		os.Exit(1)
	}

	if err := run(context.Background(), *tokenPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func LoadToken(filePath string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	tok := new(oauth2.Token)
	if err := json.Unmarshal(b, tok); err != nil {
		return nil, err
	}
	return tok, nil
}

func run(ctx context.Context, tokenPath string) error {
	tok, err := LoadToken(tokenPath)
	if err != nil {
		return err
	}
	c, err := tesla.NewClient(ctx, tok)
	if err != nil {
		return err
	}

	v, err := c.Vehicles()
	if err != nil {
		return err
	}

	for i, v := range v {
		if i > 0 {
			fmt.Println("----")
		}
		fmt.Printf("VIN: %s\n", v.Vin)
		fmt.Printf("Name: %s\n", v.DisplayName)
	}
	return nil
}
