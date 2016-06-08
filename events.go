package plugapi

// This file contains borrowed code from
// https://github.com/go-playground/webhooks/blob/v1/webhooks.go

// Event defines a Plug Event type
type Event int

// Refer to bottom of file for a list of all the events

// ProcessPayloadFunc is a common function for payload return values
type ProcessPayloadFunc func(plug *PlugDJ, payload interface{})

func (plug *PlugDJ) emitEvent(event Event, payload interface{}) {
	fn := plug.eventFuncs[event]
	if fn != nil {
		go fn(plug, payload)
	}
}

// List of Event types
const (
	AdvanceEvent                Event = iota // = "advance"
	BanEvent                                 // = "ban"
	BoothLockedEvent                         // = "boothLocked"
	ChatEvent                                // = "chat"
	ChatCommandEvent                         // = "command"
	ChatDeleteEvent                          // = "chatDelete"
	ChatLevelUpdateEvent                     // = "roomMinChatLevelUpdate"
	CommandEvent                             // = "command"
	DJListCycleEvent                         // = "djListCycle"
	DJListUpdateEvent                        // = "djListUpdate"
	DJListLockedEvent                        // = "djListLocked"
	EarnEvent                                // = "earn"
	FollowJoinEvent                          // = "followJoin"
	FloodChatEvent                           // = "floodChat"
	FriendRequestEvent                       // = "friendRequest"
	GiftedEvent                              // = "gifted"
	GrabEvent                                // = "grab"
	KillSessionEvent                         // = "killSession"
	MaintModeEvent                           // = "plugMaintenance"
	MaintModeAlertEvent                      // = "plugMaintenanceAlert"
	ModerateAddDjEvent                       // = "modAddDJ"
	ModerateAddWaitlistEvent                 // = "modAddWaitList"
	ModerateAmbassadorEvent                  // = "modAmbassador"
	ModerateBanEvent                         // = "modBan"
	ModerateMoveDjEvent                      // = "modMoveDJ"
	ModerateMuteEvent                        // = "modMute"
	ModerateRemoveDjEvent                    // = "modRemoveDJ"
	ModerateRemoveWaitlistEvent              // = "modRemoveWaitList"
	ModerateSkipEvent                        // = "modSkip"
	ModerateStaffEvent                       // = "modStaff"
	NotifyEvent                              // = "notify"
	PdjMessageEvent                          // = "pdjMessage"
	PdjUpdateEvent                           // = "pdjUpdate"
	PingEvent                                // = "ping"
	PlaylistCycleEvent                       // = "playlistCycle"
	RequestDurationEvent                     // = "requestDuration"
	RequestDurationRetryEvent                // = "requestDurationRetry"
	RoomChangeEvent                          // = "roomChanged"
	RoomDescriptionUpdateEvent               // = "roomDescriptionUpdate"
	RoomJoinEvent                            // = "roomJoin"
	RoomNameUpdateEvent                      // = "roomNameUpdate"
	RoomVoteSkipEvent                        // = "roomVoteSkip"
	RoomWelcomeUpdateEvent                   // = "roomWelcomeUpdate"
	SessionCloseEvent                        // = "sessionClose"
	SkipEvent                                // = "skip"
	StrobeToggleEvent                        // = "strobeToggle"
	UserCounterUpdateEvent                   // = "userCounterUpdate"
	UserFollowEvent                          // = "userFollow"
	UserJoinEvent                            // = "userJoin"
	UserLeaveEvent                           // = "userLeave"
	UserUpdateEvent                          // = "userUpdate"
	VoteEvent                                // = "vote"
)
