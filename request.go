package plugapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	// log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"strings"
	// "net/url"
	"net/http"
)

func (plug *PlugDJ) authenticateUser() error {
	resp, err := plug.web.Get(plug.config.BaseURL + "/")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ErrUnknownResponse{resp, "/"}
	}

	// non authenticated requests contain a little js snippet
	// beginning with: var _csrf="<60 char token"
	// we need to use this for all future requests
	var csrf string

	// we use a scanner so we don't have to read it all at once
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		text := scanner.Text()

		if pos := strings.Index(text, `_csrf="`); pos != -1 {
			csrf = text[pos+7 : pos+67] // get 60 characters, skipping the search term
			break                       // no longer need to scan now we've found the token!
		}
	}

	// scans stop on error or file finish. let's check for an error
	if err := scanner.Err(); err != nil {
		return err
	}

	// still default value, so could not find the csrf token
	if csrf == "" {
		return errors.New("dubapi: could not obtain csrf token")
	}
	plug.Log.WithField("_csrf", csrf).Info("found csrf token")

	plug.Log.Info("Attempting to log in...")

	data := map[string]string{
		"csrf":     csrf,
		"email":    plug.config.Email,
		"password": plug.config.Password,
	}

	// try to log in
	resp, err = plug.Post(AuthLoginEndpoint, data)
	if err != nil {
		return err
	}

	quickread(resp.Body)
	resp.Body.Close()

	if resp.StatusCode == 401 {
		return ErrAuthentication
	} else if resp.StatusCode != 200 {
		return ErrUnknownResponse{resp, "/"}
	}

	return nil
}

func quickread(reader io.ReadCloser) {
	b, _ := ioutil.ReadAll(reader)
	fmt.Printf("%s\n", b)
}

func (plug *PlugDJ) connectSocket() error {
	// Socket connections depend on a few things from plug:
	// - the actual socket url (_gws)
	// - the server time (_st)
	// - the ???? (_jm)
	// This works in a similar way to getting our csrf token
	// in authenticateUser above!
	// Let's go ahead and grab that!
	// plug.web.Get(url)
	return errors.New("not implemented yet!")
}

// Get information about ourselves
// func (c *apiClient) loadSession() error {
// 	endpoint := ""

// 	userdata := &User{}
// 	err := c.GetData(endpoint, userdata)
// 	if err != nil {
// 		return err
// 	}
// 	c.plug.User = userdata
// 	return nil
// }

// Get makes a connection and returns the response,
// handling any reauthentications required
// func (c *apiClient) Get(endpoint string) (*http.Response, error) {
// 	target := c.plug.config.BaseURL + endpoint

// 	resp, err := c.client.Get(target)

// 	// If we need authentication
// 	if (err != nil) && (err.(*url.Error).Err == ErrAuthenticationRequired) {
// 		// Try to authenticate
// 		if err := c.authenticate(); err != nil {
// 			return nil, err
// 		}

// 		// Resubmit the request
// 		resp, err = c.client.Get(target)

// 		// Do we still need authentication?
// 		if (err != nil) && (err.(*url.Error).Err == ErrAuthenticationRequired) {
// 			// Don't try again. Return "auth failed"
// 			return nil, ErrAuthentication
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, err
// }

// Post makes a post request with the map provided as json
func (plug *PlugDJ) Post(endpoint string, data map[string]string) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := plug.web.Post(plug.config.BaseURL+"/_"+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Responses by plugtrack are like:
// {"code":200,"message":"OK","data":{...}}
// So we need to parse the responses and pick out the data
type apiData struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// GetData allows you to receive info as a struct
// func (c *http.Client) GetData(endpoint string, v interface{}) error {
// 	resp, err := c.Get(endpoint)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	data := &apiData{}

// 	switch v.(type) {
// 	case *User:
// 		data.Data = v.(*User)
// 	case *Room:
// 		data.Data = v.(*Room)
// 	default:
// 		return ErrUnknownData
// 	}

// 	// str, err := ioutil.ReadAll(resp.Body)
// 	// fmt.Println(string(str))
// 	// if true {
// 	// 	return nil
// 	// }

// 	err = json.NewDecoder(resp.Body).Decode(data)
// 	if err != nil {
// 		return err
// 	}

// 	if data.Code != 200 {
// 		return &ErrDataRequestError{data, endpoint}
// 	}

// 	return nil
// }
