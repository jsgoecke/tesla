# tesla

<img src="https://raw.githubusercontent.com/bogosj/tesla/main/images/gotesla.png" width="200px">

This library provides a wrapper around the API to easily query and command a [Tesla](https://www.tesla.com/) remotely in Go.

## Library Documentation

[https://godoc.org/github.com/bogosj/tesla](https://godoc.org/github.com/bogosj/tesla)

## API Documentation

[View Tesla JSON API Documentation](https://tesla-api.timdorr.com/)

This is unofficial documentation of the Tesla JSON API used by the iOS and Android apps. The API provides functionality to monitor and control Telsa vehicles remotely. The project provides both a documention of the API and a Go library for accessing it.

## Installation

```
go get github.com/bogosj/tesla
```

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

## Credits

Thank you to [Tim Dorr](https://github.com/timdorr) who did the heavy lifting to document the Tesla API and also created the [model-s-api Ruby Gem](https://github.com/timdorr/model-s-api).

Thank you to [jsgoecke](https://github.com/jsgoecke) from whom this project is forked.

## Copyright & License

Copyright (c) 2016-2021 Jason Goecke.

Copyright (c) 2021-present James Bogosian.

Released under the terms of the MIT license. See LICENSE for details.
