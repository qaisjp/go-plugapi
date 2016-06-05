package plugapi

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// Signature of all action handlers
type actionHandler func(plug *PlugDJ, msg *Message)

// All messages received by the
// WS server have this structure
// TODO: Should this be exported?
type Message struct {
	Action    string      `json:"a"`
	Parameter interface{} `json:"p"`
	Time      int64       `json:"t"`
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
	// string found in Message.Action
	actions["ack"] = handleAction_ack
	actions["chat"] = handleAction_chat
}

// Base action that executes the correct handler
// or do some debug outputs if the handler
// does not exist for the given message.
func handleAction(plug *PlugDJ, msg *Message) {
	handler, ok := actions[msg.Action]
	if ok {
		// a handler exists, lets call it
		handler(plug, msg)
		return
	}

	// Default action behaviour
	plug.Log.WithFields(log.Fields{"message": msg}).Infoln("Could not handle socket message")
}

func handleAction_ack(plug *PlugDJ, msg *Message) {
	ack := plug.ack
	param, ok := msg.Parameter.(string)
	if !ok {
		ack <- errors.New("ws: 'ack' > p is not a string")
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

func handleAction_chat(plug *PlugDJ, msg *Message) {
	plug.Log.WithField("data", msg.Parameter).Infoln("Chat message...")
}
