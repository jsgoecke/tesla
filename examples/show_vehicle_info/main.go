package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/bogosj/tesla"
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

func run(ctx context.Context, tokenPath string) error {
	c, err := tesla.NewClient(ctx, tesla.WithTokenFile(tokenPath))
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
