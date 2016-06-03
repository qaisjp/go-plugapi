package plugapi

type User struct{}

// Booth is the data about the current queue
type Booth struct {
	CurrentDJ   int
	IsLocked    bool // Can users join?
	ShouldCycle bool // Will the queue progress automatically?
	WaitingDJs  []User
}

// Room contains metadata about the room
type Room struct {
	Description      string
	Favourite        bool   // Does the logged in user love this room?
	Guests           int    // ?? the below comment is correct
	Population       int    // ?? the above comment is incorrect
	HostID           int    // User ID of the room owner. default: -1
	HostName         string // Username of the room owner
	ID               int    // unique room identifier. default: -1
	MinimumChatLevel int    // power required to speak. default: 1 (POSITIVE 1)
	Name             string // name of the room
	Slug             string // string shortname
	WelcomeMessage   string // the welcome message on entering
}

// Event defines a Dubtrack Event type
type Event string

// List of Event types
const (
	AdvanceEvent                Event = "advance"
	BanEvent                    Event = "ban"
	BoothLockedEvent            Event = "boothLocked"
	ChatEvent                   Event = "chat"
	ChatCommandEvent            Event = "command"
	ChatDeleteEvent             Event = "chatDelete"
	ChatLevelUpdateEvent        Event = "roomMinChatLevelUpdate"
	CommandEvent                Event = "command"
	DJListCycleEvent            Event = "djListCycle"
	DJListUpdateEvent           Event = "djListUpdate"
	DJListLockedEvent           Event = "djListLocked"
	EarnEvent                   Event = "earn"
	FollowJoinEvent             Event = "followJoin"
	FloodChatEvent              Event = "floodChat"
	FriendRequestEvent          Event = "friendRequest"
	GiftedEvent                 Event = "gifted"
	GrabEvent                   Event = "grab"
	KillSessionEvent            Event = "killSession"
	MaintModeEvent              Event = "plugMaintenance"
	MaintModeAlertEvent         Event = "plugMaintenanceAlert"
	ModerateAddDjEvent          Event = "modAddDJ"
	ModerateAddWaitlistEvent    Event = "modAddWaitList"
	ModerateAmbassadorEvent     Event = "modAmbassador"
	ModerateBanEvent            Event = "modBan"
	ModerateMoveDjEvent         Event = "modMoveDJ"
	ModerateMuteEvent           Event = "modMute"
	ModerateRemoveDjEvent       Event = "modRemoveDJ"
	ModerateRemoveWaitlistEvent Event = "modRemoveWaitList"
	ModerateSkipEvent           Event = "modSkip"
	ModerateStaffEvent          Event = "modStaff"
	NotifyEvent                 Event = "notify"
	PdjMessageEvent             Event = "pdjMessage"
	PdjUpdateEvent              Event = "pdjUpdate"
	PingEvent                   Event = "ping"
	PlaylistCycleEvent          Event = "playlistCycle"
	RequestDurationEvent        Event = "requestDuration"
	RequestDurationRetryEvent   Event = "requestDurationRetry"
	RoomChangeEvent             Event = "roomChanged"
	RoomDescriptionUpdateEvent  Event = "roomDescriptionUpdate"
	RoomJoinEvent               Event = "roomJoin"
	RoomNameUpdateEvent         Event = "roomNameUpdate"
	RoomVoteSkipEvent           Event = "roomVoteSkip"
	RoomWelcomeUpdateEvent      Event = "roomWelcomeUpdate"
	SessionCloseEvent           Event = "sessionClose"
	SkipEvent                   Event = "skip"
	StrobeToggleEvent           Event = "strobeToggle"
	UserCounterUpdateEvent      Event = "userCounterUpdate"
	UserFollowEvent             Event = "userFollow"
	UserJoinEvent               Event = "userJoin"
	UserLeaveEvent              Event = "userLeave"
	UserUpdateEvent             Event = "userUpdate"
	VoteEvent                   Event = "vote"
)

// List of individual REST endpoints
const (
	AuthLoginEndpoint string = "/auth/login"
	RoomJoinEndpoint  string = "/rooms/join"

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

	UserInfoEndpoint       string = "/users/me"
	UserGetAvatarsEndpoint string = "/store/inventory/avatars"
	UserSetAvatarEndpoint  string = "/users/avatar"
)
