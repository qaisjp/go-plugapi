package plugapi

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	// "strconv"
	"time"
)

func (plug *PlugDJ) connectSocket() error {
	// Socket connections depend on a few things from plug:
	// - the actual socket url (_gws)
	// - the server time (_st)
	// - the passcode allowing us to authenticate
	//   with the socket server (_jm)
	// This works in the same way as getting our csrf token
	// but with different prefixes

	// Let's go ahead and grab that!
	resp, err := plug.web.Get(plug.config.BaseURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	variables, err := scanPrefixes(resp.Body, ",_gws", ",_jm", ",_st")
	if err != nil {
		return err
	}

	// We'll use this later
	var ok bool
	plug.authCode, ok = variables[",_jm"]
	if !ok {
		panic("impossible situation - check variables map key")
	}

	// We don't want to override any forced socket urls...
	if plug.config.SocketURL == "" {
		plug.config.SocketURL, ok = variables[",_gws"]
		if !ok {
			panic("impossible situation - check variables map key")
		}
	}

	// Now we want to also store a time offset because
	// one of us (bot or plug) are ahead of each other
	timeStr, ok := variables[",_st"]
	if !ok {
		panic("impossible situation - check variables map key")
	}

	// Format of the time used (manually composed)
	// see: https://golang.org/src/time/format.go
	format := "2006-01-02 15:04:05.000000"
	theirTime, err := time.Parse(format, timeStr)
	if err != nil {
		plug.Log.WithField("_st", variables[",_st"]).Warnln("could not parse correctly")
		return errors.New("plugapi: could not parse _st correctly")
	}

	// offset the time, store it in seconds
	// note: Seconds() returns a float, and int() truncates it
	offset := int(time.Now().Sub(theirTime).Seconds())
	plug.Log.WithField("offset", offset).Debugln("Received time offset from plug.dj server")
	// fmt.Println(time.Now(), theirTime)

	// plugdj runs in their own timezone... valve time
	// now we can use this Location to handle times
	// correctly everywhere on our bot
	plug.location = time.FixedZone("plugdj", offset)

	// make a header with our origin...
	header := make(http.Header)
	header.Set("Origin", plug.config.BaseURL)

	// try to dial a connection to the websocket
	plug.Log.Debugln("Dialing websocket...")
	wss, _, err := websocket.DefaultDialer.Dial(plug.config.SocketURL, header)
	if err != nil {
		plug.Log.WithFields(log.Fields{
			"socketURL": plug.config.SocketURL,
			"baseURL":   plug.config.BaseURL,
		}).Fatalf("websocket.Dial encountered error>> %s", err)
		return err
	}

	// add the websocket to the plug obj
	plug.wss = wss

	// start listening
	plug.ack = make(chan error)
	go plug.listen()

	// Now we try to authenticate with our auth code...
	plug.Log.Debugln("Authenticating with our websocket...")
	err = plug.sendSocketJSON("auth", plug.authCode)
	if err != nil {
		plug.Log.Warnf("Failed to authenticate with our websocket >> %s\n", err)
		return err
	}

	select {
	// wait until we have successfully authenticated
	case err, failed := <-plug.ack:
		if failed {
			return err
		}
		return nil
	// or five seconds have passed
	case <-time.After(time.Second * 5):
		return errors.New("could not authenticate with WS server")
	}
}

func (plug *PlugDJ) sendSocketJSON(action string, data interface{}) error {
	body := socketMessage{
		Action:    action,
		Parameter: data,
		Time:      time.Now().In(plug.location).Unix(), // NOTE: NEEDS TO BE NUMBER NOT STRING
	}

	// plug.Log.WithField("body", body).Debugln("Sending WS data")
	return plug.wss.WriteJSON(body)
}

func (plug *PlugDJ) listen() {
	defer plug.wss.Close()
	defer close(plug.closer)
	for {
		_, data, err := plug.wss.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				plug.Log.Errorln("socket read error:", err)
				// TODO: Should we attempt to re-establish a connection?
			}
			return
		}

		// ignore messages with just "h"
		if (len(data) == 1) && (data[0] == 'h') {
			continue
		}

		// for some reason the server may send multiple messages
		var messages []json.RawMessage
		if err := json.Unmarshal([]byte(data), &messages); err != nil {
			plug.Log.WithField("data", string(data)).Warnf("ws: could not unmarshal socket array>> %s\n", err)
			return
		}

		for _, buf := range messages {
			// plug.Log.Debugln(string(buf))
			go func(buf json.RawMessage) {

				// init a message with our json.RawMessage
				// Param so that we can read it later
				msg := socketMessage{
					Parameter: new(json.RawMessage),
				}

				// unmarshal it
				if err := json.Unmarshal(buf, &msg); err != nil {
					plug.Log.WithField("data", string(data)).Warnf("ws: could not unmarshal>> %s\n", err)
					return
				}

				// do some "depointering" so that we
				// don't have to do it when handling it
				msg.Parameter = *msg.Parameter.(*json.RawMessage)

				// send it off to our socket message handler
				handleAction(plug, msg)
			}(buf)
		}

	}
}
