package plugapi

import (
	"fmt"
	"io"
	"io/ioutil"
)

type User struct{}

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

// HistoryItem is an individual item in the room history
type HistoryItem struct {
	ID    string `json:"id"`
	Media Media  `json:"media"`
	Room  struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"room"`
	Score struct {
		Grabs     int `json:"grabs"`
		Listeners int `json:"listeners"`
		Negative  int `json:"negative"`
		Positive  int `json:"positive"`
		Skipped   int `json:"skipped"`
	} `json:"score"`
	Timestamp string `json:"timestamp"` // Format: xx-xx-xx
	User      struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
}

// Room contains metadata about the room
type Room struct {
	Booth Booth `json:"booth"`
	// FX interface{} `json:"fx"`
	// Grabs interface{} `json:"grabs"`
	Meta struct {
		Description      string `json:"description"`
		Favorite         bool   `json:"favorite"`       // Does the logged in user love this room?
		Guests           int    `json:"guests"`         // Number of guests connected
		HostID           int    `json:"hostID"`         // User ID of the room owner. default: -1
		HostName         string `json:"hostName"`       // Username of the room owner
		ID               int    `json:"id"`             // unique room identifier. default: -1
		MinimumChatLevel int    `json:"minChatLevel"`   // power required to speak. default: 1 (POSITIVE 1)
		Name             string `json:"name"`           // name of the room
		Population       int    `json:"population"`     // Number of real users in the room (guests excluded)
		Slug             string `json:"slug"`           // string shortname
		WelcomeMessage   string `json:"welcomemessage"` // the welcome message on entering
	} `json:"meta"`
	// Mutes interface{} `json:"mutes"`
	Playback Playback `json:"playback"`
	Role     int      `json:"role"` // OUR ROLE IN THE ROOM << DO NOT USE
	// Users interface{} `json:"users"`
	// Votes interface{} `json:"votes"`
}

// Event defines a Plug Event type
// type Event string

// List of Event types
const (
	AdvanceEvent                int = iota // = "advance"
	BanEvent                               // = "ban"
	BoothLockedEvent                       // = "boothLocked"
	ChatEvent                              // = "chat"
	ChatCommandEvent                       // = "command"
	ChatDeleteEvent                        // = "chatDelete"
	ChatLevelUpdateEvent                   // = "roomMinChatLevelUpdate"
	CommandEvent                           // = "command"
	DJListCycleEvent                       // = "djListCycle"
	DJListUpdateEvent                      // = "djListUpdate"
	DJListLockedEvent                      // = "djListLocked"
	EarnEvent                              // = "earn"
	FollowJoinEvent                        // = "followJoin"
	FloodChatEvent                         // = "floodChat"
	FriendRequestEvent                     // = "friendRequest"
	GiftedEvent                            // = "gifted"
	GrabEvent                              // = "grab"
	KillSessionEvent                       // = "killSession"
	MaintModeEvent                         // = "plugMaintenance"
	MaintModeAlertEvent                    // = "plugMaintenanceAlert"
	ModerateAddDjEvent                     // = "modAddDJ"
	ModerateAddWaitlistEvent               // = "modAddWaitList"
	ModerateAmbassadorEvent                // = "modAmbassador"
	ModerateBanEvent                       // = "modBan"
	ModerateMoveDjEvent                    // = "modMoveDJ"
	ModerateMuteEvent                      // = "modMute"
	ModerateRemoveDjEvent                  // = "modRemoveDJ"
	ModerateRemoveWaitlistEvent            // = "modRemoveWaitList"
	ModerateSkipEvent                      // = "modSkip"
	ModerateStaffEvent                     // = "modStaff"
	NotifyEvent                            // = "notify"
	PdjMessageEvent                        // = "pdjMessage"
	PdjUpdateEvent                         // = "pdjUpdate"
	PingEvent                              // = "ping"
	PlaylistCycleEvent                     // = "playlistCycle"
	RequestDurationEvent                   // = "requestDuration"
	RequestDurationRetryEvent              // = "requestDurationRetry"
	RoomChangeEvent                        // = "roomChanged"
	RoomDescriptionUpdateEvent             // = "roomDescriptionUpdate"
	RoomJoinEvent                          // = "roomJoin"
	RoomNameUpdateEvent                    // = "roomNameUpdate"
	RoomVoteSkipEvent                      // = "roomVoteSkip"
	RoomWelcomeUpdateEvent                 // = "roomWelcomeUpdate"
	SessionCloseEvent                      // = "sessionClose"
	SkipEvent                              // = "skip"
	StrobeToggleEvent                      // = "strobeToggle"
	UserCounterUpdateEvent                 // = "userCounterUpdate"
	UserFollowEvent                        // = "userFollow"
	UserJoinEvent                          // = "userJoin"
	UserLeaveEvent                         // = "userLeave"
	UserUpdateEvent                        // = "userUpdate"
	VoteEvent                              // = "vote"
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
