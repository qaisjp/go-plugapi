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
	Users    []*User   `json:"users"`
	// Votes interface{} `json:"votes"`
}

func gatherUsers(p *PlugDJ, users []int) []*User {
	results := make([]*User, 0, len(users))
	for _, uid := range users {
		for _, user := range p.Room.Users {
			if user.ID == uid {
				results = append(results, user)
			}
		}
	}
	return results
}

func (p *PlugDJ) getDJ() *User {
	r := p.Room

	if r.Booth.CurrentDJ <= 0 {
		return nil
	}

	// user := r.GetUser(r.Booth.CurrentDJ)
	// TODO: What does cacheUser do here?
	// (see room.js)

	return p.getUser(r.Booth.CurrentDJ)

}

func (p *PlugDJ) getUser(id int) *User {
	r := p.Room

	// Base case: is it ourself?
	if id == p.User.ID {
		return p.User
	}

	// Linear search for the user
	for _, user := range r.Users {
		if user.ID == id {
			return user
		}
	}

	// Couldn't find it, sorry.
	return nil
}

func (p *PlugDJ) getDJs() []*User {
	return gatherUsers(p, p.Room.Booth.WaitingDJs)
}

func (p *PlugDJ) removeUser(id int) (u *User) {
	index := -1
	for i, user := range p.Room.Users {
		if user.ID == id {
			index = i
			u = user
			break
		}
	}

	if index == -1 {
		return nil
	}

	p.Room.Users = append(p.Room.Users[:index], p.Room.Users[index+1:]...)
	return
}

func (p *PlugDJ) addUser(u User) {
	p.removeUser(u.ID)
	p.Room.Users = append(p.Room.Users, &u)
}
