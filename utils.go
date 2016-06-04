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
