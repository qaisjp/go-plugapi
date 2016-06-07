package plugapi

type ChatPayload struct {
	Message    string `json:"message"`
	UserID     int    `json:"uid"`
	Username   string `json:"un"`
	ChatID     int    `json:"cid"`
	Subscriber int    `json:"sub"` // subscriber state. usually 0/1 for no/yes
}
