# tesla

<img src="https://raw.githubusercontent.com/bogosj/tesla/main/.github/images/gotesla.png" width="200px">

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

Examples can be found in the [/examples project directory](examples).

### OAuth Token

One way to acquire an OAuth token is to run cmd/login.

```sh

# cd cmd/login

# go run . -o ~/tesla.token
✔ Username: email@gmail.com
✔ Password: ***
Passcode: 463932
```

This will output a token to the `tesla.token` file in your home directory.

## Differences from jsgoecke/tesla

### Streaming API

The implementation of the Streaming API is not working. The code was removed [in this commit](https://github.com/bogosj/tesla/commit/19f79e1dc7a6c5d5ea5d5c8e0f4f0f2c42673404). If you are interested in getting this working again feel free to send a PR reverting these changes and providing a fix.

## Credits

Thank you to [Tim Dorr](https://github.com/timdorr) who did the heavy lifting to document the Tesla API and also created the [model-s-api Ruby Gem](https://github.com/timdorr/model-s-api).

Thank you to [jsgoecke](https://github.com/jsgoecke) from whom this project is forked.

## Copyright & License

Copyright (c) 2016-2021 Jason Goecke.

Copyright (c) 2021-present James Bogosian.

Released under the terms of the MIT license. See LICENSE for details.
