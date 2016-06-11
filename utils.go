package plugapi

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type User struct {
	ID       int
	Role     int
	Username string
}

// Booth is the data about the current queue
type Booth struct {
	CurrentDJ   int   `json:"currentDJ"`
	IsLocked    bool  `json:"isLocked"`    // Can users join?
	ShouldCycle bool  `json:"shouldCycle"` // Will the queue progress automatically?
	WaitingDJs  []int `json:"waitingDJs"`
}

type Media struct {
	Author   string `json:"author"`
	CID      string `json:"cid"`
	Duration int    `json:"duration"` // default: -1
	Format   int    `json:"format"`   // default: -1
	ID       int    `json:"id"`       // default: -1
	Image    string `json:"image"`
	Title    string `json:"title"`
}

// Playback metadata about an existing play (note, not the song)
type Playback struct {
	HistoryID  string `json:"historyID"`
	Media      Media  `json:"media"`
	PlaylistID int    `json:"playlistID"` // default: -1
	StartTime  string `json:"startTime"`
}

// PlayScore is the score of an individual song
type PlayScore struct {
	Grabs     int `json:"grabs"`
	Listeners int `json:"listeners"`
	Negative  int `json:"negative"`
	Positive  int `json:"positive"`
	Skipped   int `json:"skipped"`
}

// HistoryItem is an individual item in the room history
type HistoryItem struct {
	ID    string `json:"id"`
	Media Media  `json:"media"`
	Room  struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"room"`
	Score     PlayScore `json:"score"`
	Timestamp string    `json:"timestamp"` // Format: xx-xx-xx
	User      struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
}

type chatMessageType int

const (
	RegularChatMessage chatMessageType = iota
	EmoteChatMessage
)

// List of individual REST endpoints
const (
	AuthLoginEndpoint string = "/auth/login"

	ChatDeleteEndpoint string = "/chat/"
	HistoryEndpoint    string = "/rooms/history"
	PlaylistEndpoint   string = "/playlists"

	ModerateAddDJEndpoint       string = "/booth/add"
	ModerateBanEndpoint         string = "/bans/add"
	ModerateBoothEndpoint       string = "/booth"
	ModerateMoveDJEndpoint      string = "/booth/move"
	ModerateMuteEndpoint        string = "/mutes"
	ModeratePermissionsEndpoint string = "/staff/update"
	ModerateRemoveDJEndpoint    string = "/booth/remove/"
	ModerateSkipEndpoint        string = "/booth/skip"
	ModerateStaffEndpoint       string = "/staff/"
	ModerateUnbanEndpoint       string = "/bans/"
	ModerateUnmuteEndpoint      string = "/mutes/"

	SkipMeEndpoint         string = "/booth/skip/me"
	RoomCycleBoothEndpoint string = "/booth/cycle"
	RoomLockBoothEndpoint  string = "/booth/lock"
	RoomInfoEndpoint       string = "/rooms/update"
	RoomJoinEndpoint       string = "/rooms/join"
	RoomStateEndpoint      string = "/rooms/state"

	UserInfoEndpoint       string = "/users/me"
	UserGetAvatarsEndpoint string = "/store/inventory/avatars"
	UserSetAvatarEndpoint  string = "/users/avatar"
)

// TODO: quickread is a debug function
// to print and return all of an
// io.Reader's contents
func quickread(reader io.Reader) []byte {
	b, _ := ioutil.ReadAll(reader)
	fmt.Printf("%s\n", b)
	return b
}

type IntBool bool

func (b *IntBool) UnmarshalJSON(data []byte) error {
	n, err := strconv.Atoi(string(data))

	if err != nil {
		return err
	}

	if n != 0 && n != 1 {
		return errors.New("IntBool: got non 0 or 1 value")
	}

	*b = n == 1
	return nil
}
