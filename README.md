# go-dubapi
[![GoDoc](https://godoc.org/github.com/qaisjp/go-dubapi?status.svg)](https://godoc.org/github.com/qaisjp/go-dubapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/qaisjp/go-dubapi)](https://goreportcard.com/report/github.com/qaisjp/go-dubapi)

go-dubapi is a Go client library for the Dubtrack API.

This library is being developed for a bot I am running at [just-a-chill-room](http://just-a-chill-room.net), so API methods will likely be implemented in the order that they are needed by that application. You can track the status of implementation in [this Google spreadsheet][roadmap]. Eventually, I would like to cover the entire Dubtrack API, so contributions are of course always welcome.

[roadmap]: https://docs.google.com/spreadsheets/d/1q6Qj4G9yxQ2_L2pmw2LE5Gpu7ndPctllccB8T5diR_I/edit?usp=sharing

The calling pattern and structure or this library is **subject to change at any time**, so all suggestions for API improvements are welcome and encouraged.

Projects currently being referred to for style (this is my first time writing a framework in Go) and implementation:
* https://github.com/anjanms/DubAPI (implementation)
* https://github.com/huandu/facebook (style)
* https://github.com/shkh/lastfm-go (style)
* https://github.com/go-playground/webhooks (event implementation and style)
