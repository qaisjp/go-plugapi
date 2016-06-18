package plugapi

import "sync"

// Room contains metadata about the room
// TODO: Unexport this.
type Room struct {
	sync.RWMutex
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
	Playback *Playback `json:"playback"`
	Role     int       `json:"role"` // OUR ROLE IN THE ROOM << DO NOT USE
	users    []User    `json:"users"`
	// Votes interface{} `json:"votes"`
}

func gatherUsers(r *Room, users []int) []User {
	r.RLock()
	defer r.RUnlock()

	results := make([]User, 0, len(users))
	for _, uid := range users {
		for _, user := range r.users {
			if user.ID == uid {
				results = append(results, user)
			}
		}
	}
	return results
}

func (r *Room) getDJ() *User {
	r.RLock()

	if r.Booth.CurrentDJ <= 0 {
		return nil
	}

	// user := r.GetUser(r.Booth.CurrentDJ)
	// TODO: What does cacheUser do here?
	// (see room.js)

	dj := r.Booth.CurrentDJ

	r.RUnlock()
	return r.getUser(dj)

}

func (r *Room) getUser(id int) *User {

	// Base case: is it ourself?
	// TODO: Needs reference to self
	// if id == p.User.ID {
	// return p.User
	// }

	r.RLock()
	defer r.RUnlock()

	// Linear search for the user
	for _, user := range r.users {
		if user.ID == id {
			return &user
		}
	}

	// Couldn't find it, sorry.
	return nil
}

func (r *Room) getDJs() []User {
	return gatherUsers(r, r.Booth.WaitingDJs)
}

func (r *Room) removeUser(id int) (u *User) {
	r.Lock()
	defer r.Unlock()

	index := -1
	for i, user := range r.users {
		if user.ID == id {
			index = i
			u = &user
			break
		}
	}

	if index == -1 {
		return nil
	}

	r.users = append(r.users[:index], r.users[index+1:]...)
	return
}

func (r *Room) addUser(u User) {
	r.Lock()
	defer r.Unlock()

	r.removeUser(u.ID)
	r.users = append(r.users, u)
}

func (r *Room) GetUsers() (u []User) {
	r.RLock()
	copy(u, r.users)
	r.RUnlock()
	return
}
