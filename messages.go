package plugapi

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// Signature of all action handlers
type actionHandler func(plug *PlugDJ, msg json.RawMessage)

// All messages received by the
// WS server have this structure
type socketMessage struct {
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
	// string found in socketMessage.Action
	actions["ack"] = handleAction_ack
	actions["chat"] = handleAction_chat
}

// Base action that executes the correct handler
// or do some debug outputs if the handler
// does not exist for the given message.
func handleAction(plug *PlugDJ, msg socketMessage) {
	handler, ok := actions[msg.Action]
	if ok {
		// a handler exists, lets call it,
		// but give it the param directly
		handler(plug, msg.Parameter.(json.RawMessage))
		return
	}

	// Default action behaviour
	msg.Parameter = string(msg.Parameter.(json.RawMessage))
	plug.Log.WithFields(log.Fields{"message": msg}).Debugln("WS: ??:")
}

func handleAction_ack(plug *PlugDJ, msg json.RawMessage) {
	ack := plug.ack

	var param string
	err := json.Unmarshal(msg, &param)
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

func handleAction_chat(plug *PlugDJ, msg json.RawMessage) {
	var payload ChatPayload
	json.Unmarshal(msg, &payload)
	plug.emitEvent(ChatEvent, payload)
}
