package plugapi

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// Signature of all action handlers
type actionHandler func(plug *PlugDJ, msg *messageOut)

// All messages received by the
// WS server have this structure
type messageOut struct {
	Action    string      `json:"a"`
	Parameter interface{} `json:"p"`
	Time      int64       `json:"t"`
}

type MessageIn struct {
	Action    string          `json:"a"`
	Parameter json.RawMessage `json:"p"`
	Time      int64           `json:"t"`
}

// Map of all actions we can handle
// - map key is the "Action" string
// - value is a function to handle the action
var actions map[string]actionHandler

func init() {
	actions = make(map[string]actionHandler)

	// add all your new actions here make
	// sure your function is with the form:
	// 		handleAction_ACTIONNAME
	// where ACTIONNAME is the exact
	// string found in messageOut.Action
	actions["ack"] = handleAction_ack
	actions["chat"] = handleAction_chat
}

// Base action that executes the correct handler
// or do some debug outputs if the handler
// does not exist for the given message.
func handleAction(plug *PlugDJ, msg *messageOut) {
	handler, ok := actions[msg.Action]
	if ok {
		// a handler exists, lets call it
		handler(plug, msg)
		return
	}

	// Default action behaviour
	plug.Log.WithFields(log.Fields{"message": msg}).Debugln("Could not handle socket message")
}

func handleAction_ack(plug *PlugDJ, msg *messageOut) {
	ack := plug.ack

	var param string
	err := json.Unmarshal(msg.Parameter.([]byte), &param)
	if err != nil {
		ack <- err
		return
	}

	if param, err := strconv.Atoi(string(param)); err != nil {
		plug.Log.WithField("error", err).Warnln("could not read 'ack' param value")
		ack <- errors.New("ws: 'ack' > Parameter not integer")
	} else if param == 1 {
		close(ack)
	} else {
		plug.Log.WithField("Param", param).Warnln("Parameter is not equal 1")
		ack <- errors.New("ws: ack > Parameter not 1")
	}
}

func handleAction_chat(plug *PlugDJ, msg *messageOut) {
	var payload ChatPayload
	json.Unmarshal(msg.Parameter.([]byte), &payload)
	plug.emitEvent(ChatEvent, payload)
}
