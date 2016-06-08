package plugapi

// This file is dedicated towards structs that definitely need to be
// accessed by other packages.

// For pauloads, intermediate structs should be used to handle incoming json
// unless it can be directly unmarshalled (with the exception of IntBool)

type ChatPayload struct {
	Message string // The chat message
	User    *User  // Who it came from
	Type    chatMessageType
}

type AdvancePayload struct {
	CurrentDJ *User `json:"c"` // TODO: Write unmarshaler for User, with reference to original plug obj??
	DJs       []*User
	LastPlay  *struct {
		DJ    *User
		Media Media
		Score PlayScore
	}
	Playback *Playback
}

// CommandHandler is the structure of a function that can handle a chat command
type CommandHandler func(data CommandData, command string, args ...string)

// CommandData is used when receiving individual command events
type CommandData struct {
	Plug      *PlugDJ
	User      *User
	MessageID string
}
