# go-plugapi
[![GoDoc](https://godoc.org/github.com/qaisjp/go-plugapi?status.svg)](https://godoc.org/github.com/qaisjp/go-plugapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/qaisjp/go-plugapi)](https://goreportcard.com/report/github.com/qaisjp/go-plugapi)

go-plugapi is a Go client library for the Plug.dj API.

This library is being developed for a bot I am running at [just-a-chill-room](http://just-a-chill-room.net), so API methods will likely be implemented in the order that they are needed by that application. You can track the status of implementation in [this Google spreadsheet][roadmap] (not ready yet). Eventually, I would like to cover the entire Plug.dj API, so contributions are of course always welcome.

[roadmap]: #

The calling pattern and structure or this library is **subject to change at any time**, so all suggestions for API improvements are welcome and encouraged.

Projects currently being referred to for style (this is my first time writing a framework in Go) and implementation:
* https://github.com/plugCubed/plugAPI (implementation)
* https://github.com/huandu/facebook (style)
* https://github.com/shkh/lastfm-go (style)
* https://github.com/go-playground/webhooks (event implementation and style)
