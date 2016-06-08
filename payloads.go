package plugapi

// This file is dedicated towards payloads for use in events.
// Intermediate structs should be used to handle incoming json
// unless it can be directly unmarshalled (with the exception of IntBool)

type ChatPayload struct {
	Message    string  `json:"message"`
	UserID     int     `json:"uid"`
	Username   string  `json:"un"`
	ChatID     int     `json:"cid"`
	Subscriber IntBool `json:"sub"` // subscriber state. usually 0/1 for no/yes
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
