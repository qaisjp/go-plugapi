package plugapi

import (
	"golang.org/x/net/publicsuffix"
	// "io/ioutil"
	// "crypto/sha512"
	"errors"
	// "encoding/hex"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	// "strings"
)

// PlugDJ is the individual user connected to plug
type PlugDJ struct {
	config *Config
	Room   *Room
	User   *User
	Log    *log.Logger

	web                 *http.Client
	authCode            string
	currentlyConnecting bool
}

// Config is the configuration for logging into plug
type Config struct {
	Email    string
	Password string
	BaseURL  string
	Log      *log.Logger
}

// New returns an authenticated User
func New(config Config) (*PlugDJ, error) {
	if config.Log == nil {
		config.Log = log.New()
	}

	// make sure they gave us a valid email address
	if config.Email == "" {
		return nil, ErrAuthenticationRequired
	}

	// default base url
	if config.BaseURL == "" {
		config.BaseURL = "https://plug.dj"
	}

	// Double check the url...
	if _, err := url.Parse(config.BaseURL); err != nil {
		return nil, errors.New("plugapi: invalid url provided")
	}

	plug := &PlugDJ{config: &config, Log: config.Log}

	// was this just used to uniquely get a fucking jar?!
	// hash := sha512.Sum512([]byte(config.Email + config.Password))
	// cookieHash := hex.EncodeToString(hash[:])
	// plug.Log.WithField("cookieHash", cookieHash).Info("from new")

	// create a cookie jar to make sure we can do further requests
	opts := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	cookieJar, _ := cookiejar.New(&opts)

	// create our web client so that we can make REST requests
	plug.web = &http.Client{Jar: cookieJar}

	if err := plug.authenticateUser(); err != nil {
		return nil, err
	}

	plug.Log.Info("Running go-plugapi")
	return plug, nil
}

func (plug *PlugDJ) JoinRoom(slug string) error {
	// prevent multiple simultaneous connections
	if plug.currentlyConnecting {
		return errors.New("plugapi: already connecting to a room")
	}

	// get our list of cookies so that we can check if we are logged in
	url, _ := url.Parse(plug.config.BaseURL) // we ignore parse errors...
	cookies := plug.web.Jar.Cookies(url)
	sessionSet := false

	// iterate through and check for a "session" cookie
	for _, cookie := range cookies {
		sessionSet = cookie.Name == "session"
		if sessionSet {
			break
		}
	}

	// tell them they are not logged in
	if !sessionSet {
		// plugapi waits "a frame" to try again, we won't do this.
		return errors.New("plugapi: not logged in")
	}

	plug.Room = &Room{Slug: slug}
	plug.currentlyConnecting = true

	// NOTE: Reference > queueConnectSocket(roomSlug) < is now called
	// This tells the queue to call > connectSocket(roomSlug) <

	// Now we need to make a socket connection
	if err := plug.connectSocket(); err != nil {
		return err
	}

	// TODO: Should this be queued?
	resp, err := plug.Post(RoomJoinEndpoint, map[string]string{"slug": slug})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return errors.New("plugapi: invalid room url")
	}

	quickread(resp.Body)

	return nil
}

func (plug *PlugDJ) GetRoomData(slug string) (*Room, error) {
	room := &Room{}
	// url := strings.Replace(roomEndpoint, "{room_slug}", slug, 1)
	// err := plug.Client.GetData(url, room)

	// if err != nil {
	// 	return nil, err
	// }

	return room, nil
}
