# tesla
[![wercker status](https://app.wercker.com/status/c8e21c53ed5763b0b58f763670753732/m "wercker status")](https://app.wercker.com/project/bykey/c8e21c53ed5763b0b58f763670753732)

<p align="center">
  <img src="https://raw.githubusercontent.com/jsgoecke/tesla/master/images/gotesla.png">
  <img src="https://raw.githubusercontent.com/jsgoecke/tesla/master/images/tesla-model-s.png">
  <img src="https://raw.githubusercontent.com/jsgoecke/tesla/master/images/tesla-model-x.png">
</p>

This library provides a wrapper around the API to easily query and command the a [Tesla Model S](https://www.teslamotors.com/models) remotely in Go.

## Library Documentation

[https://godoc.org/github.com/jsgoecke/tesla](https://godoc.org/github.com/jsgoecke/tesla)

## API Documentation

[View Tesla JSON API Documentation](https://tesla-api.timdorr.com/)

This is unofficial documentation of the Tesla JSON API used by the iOS and Android apps. The API provides functionality to monitor and control the Model S (and future Tesla vehicles) remotely. The project provides both a documention of the API and a Go library for accessing it.

## Installation

```
go get github.com/jsgoecke/tesla
```

## Tokens

You may get your tokens to use as client_id and client_secret [here](http://pastebin.com/fX6ejAHd).

## Usage

Here's an example (more in the /examples project directory):

```go
func main() {
	client, err := tesla.NewClient(
		&tesla.Auth{
			ClientID:     os.Getenv("TESLA_CLIENT_ID"),
			ClientSecret: os.Getenv("TESLA_CLIENT_SECRET"),
			Email:        os.Getenv("TESLA_USERNAME"),
			Password:     os.Getenv("TESLA_PASSWORD"),
		})
	if err != nil {
		panic(err)
	}

	vehicles, err := client.Vehicles()
	if err != nil {
		panic(err)
	}

	vehicle := vehicles[0]
	status, err := vehicle.MobileEnabled()
	if err != nil {
		panic(err)
	}

	fmt.Println(status)
	fmt.Println(vehicle.HonkHorn())

	// Autopark
	// Use with care, as this will move your car
	vehicle.AutoparkForward()
	vehicle.AutoparkReverse()
	// Use with care, as this will move your car

	// Stream vehicle events
	eventChan, errChan, err := vehicle.Stream()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		for {
			select {
			case event := <-eventChan:
				eventJSON, _ := json.Marshal(event)
				fmt.Println(string(eventJSON))
			case err = <-errChan:
				fmt.Println(err)
				if err.Error() == "HTTP stream closed" {
					fmt.Println("Reconnecting!")
					eventChan, errChan, err := vehicle.Stream()
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}
```

## Examples

* [Commanding a Tesla Model S with the Amazon Echo](https://medium.com/@jsgoecke/commanding-a-tesla-model-s-with-the-amazon-echo-a06f975364b8#.xoctd3yni)

## Pull Requests

I appreciate all the pull requests to date, and have merged them when I can. I would kindly ask going forward if you do send a PR, please ensure it has the corresponding unit test, or unit test change, that passes. I have tried to write tests after the fact, but that is not best practice. Thank you!

## Credits

Thank you to [Tim Dorr](https://github.com/timdorr) who did the heavy lifting to document the Tesla API and also created the [model-s-api Ruby Gem](https://github.com/timdorr/model-s-api).

## Copyright & License

Copyright (c) 2016-Present Jason Goecke. Released under the terms of the MIT license. See LICENSE for details.
