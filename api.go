package plugapi

import (
	"golang.org/x/net/publicsuffix"
	// "io/ioutil"
	// "crypto/sha512"
	"errors"
	// "encoding/hex"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	// "strings"
)

// PlugDJ is the individual user connected to plug
type PlugDJ struct {
	config *Config
	Room   *Room
	User   *User
	Log    *log.Logger

	web                 *http.Client
	wss                 *websocket.Conn
	authCode            string
	currentlyConnecting bool
	location            *time.Location

	closer chan struct{}

	// only used once in sockets.go to determine
	// whether WS server authentication succeeded
	ack chan error
}

// Config is the configuration for logging into plug
type Config struct {
	Email     string
	Password  string
	BaseURL   string
	SocketURL string
	Log       *log.Logger
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

	// a closer so that we can close any goroutines we have created
	plug.closer = make(chan struct{})

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

func (plug *PlugDJ) Close() {
	plug.Log.Debugln("plugapi will now close")

	if plug.wss != nil {
		// To cleanly close a connection, a client should send a close
		// frame and wait for the server to close the connection.
		err := plug.wss.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			plug.Log.Println("write close:", err)
			return
		}

		select {
		// Our socket ticker will receive an error
		// the ticker will write to plug.closer when the
		// error is a "CloseNormalClosure" error
		case <-plug.closer:
			plug.Log.Debugln("sockets closed successfully")
		// As a backup, we wait a second instead.
		case <-time.After(time.Second):
			plug.Log.Warnln("sockets took too long to close")
		}

		// Now we close our clientside connection
		plug.wss.Close()
	}
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

	plug.currentlyConnecting = true

	// NOTE: Reference > queueConnectSocket(roomSlug) < is now called
	// This tells the queue to call > connectSocket(roomSlug) <

	// Now we need to make a socket connection
	if err := plug.connectSocket(); err != nil {
		return err
	}

	// TODO: Should this be queued?
	plug.Log.Debugln("Joining room...")
	resp, err := plug.Post(RoomJoinEndpoint, map[string]string{"slug": slug})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return errors.New("plugapi: invalid room url")
	}

	// Now we need to load ALL information about our current room state
	var data []*Room
	_, err = plug.GetData(RoomStateEndpoint, &data)
	if err != nil {
		return err
	}

	plug.Room = data[0]

	// See initRoom(data, callback)

	// Now we need to emit an AdvanceEvent
	// ..

	// Retrieve our history
	// resp, err = plug.Get(HistoryEndpoint)
	// if err != nil {
	// 	return err
	// }
	// quickread(resp.Body)

	// Emit something else...

	plug.currentlyConnecting = false

	return nil
}
